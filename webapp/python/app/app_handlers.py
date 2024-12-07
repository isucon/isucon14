from http import HTTPStatus
from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException, Response
from pydantic import BaseModel
from sqlalchemy import text
from sqlalchemy.engine import Connection
from ulid import ULID

from .middlewares import app_auth_middleware
from .models import (
    Chair,
    ChairLocation,
    Owner,
    PaymentToken,
    Ride,
    RideStatus,
    User,
)
from .payment_gateway import (
    PaymentGatewayPostPaymentRequest,
    UpstreamError,
    request_payment_gateway_post_payment,
)
from .sql import engine
from .utils import (
    FARE_PER_DISTANCE,
    INITIAL_FARE,
    calculate_distance,
    calculate_fare,
    secure_random_str,
    timestamp_millis,
)

router = APIRouter(prefix="/api/app")


class AppPostUsersRequest(BaseModel):
    username: str
    firstname: str
    lastname: str
    date_of_birth: str
    invitation_code: str | None = None


class AppPostUsersResponse(BaseModel):
    id: str
    invitation_code: str


@router.post("/users", status_code=HTTPStatus.CREATED)
def app_post_users(
    req: AppPostUsersRequest, response: Response
) -> AppPostUsersResponse:
    user_id = str(ULID())
    access_token = secure_random_str(32)
    invitation_code = secure_random_str(15)

    with engine.begin() as conn:
        conn.execute(
            text(
                "INSERT INTO users (id, username, firstname, lastname, date_of_birth, access_token, invitation_code) VALUES (:id, :username, :firstname, :lastname, :date_of_birth, :access_token, :invitation_code)"
            ),
            {
                "id": user_id,
                "username": req.username,
                "firstname": req.firstname,
                "lastname": req.lastname,
                "date_of_birth": req.date_of_birth,
                "access_token": access_token,
                "invitation_code": invitation_code,
            },
        )

        # 初回登録キャンペーンのクーポンを付与
        conn.execute(
            text(
                "INSERT INTO coupons (user_id, code, discount) VALUES (:user_id, :code, :discount)"
            ),
            {"user_id": user_id, "code": "CP_NEW2024", "discount": 3000},
        )

        # 招待コードを使った登録
        if req.invitation_code:
            # 招待する側の招待数をチェック
            coupons = conn.execute(
                text("SELECT * FROM coupons WHERE code = :code FOR UPDATE"),
                {"code": "INV_" + req.invitation_code},
            ).fetchall()

            if len(coupons) >= 3:
                raise HTTPException(
                    status_code=HTTPStatus.BAD_REQUEST,
                    detail="この招待コードは使用できません",
                )

            # ユーザーチェック
            inviter = conn.execute(
                text("SELECT * FROM users WHERE invitation_code = :invitation_code"),
                {"invitation_code": req.invitation_code},
            ).fetchone()

            if inviter is None:
                raise HTTPException(
                    status_code=HTTPStatus.BAD_REQUEST,
                    detail="この招待コードは使用できません。",
                )

            # 招待クーポン付与
            conn.execute(
                text(
                    "INSERT INTO coupons (user_id, code, discount) VALUES (:user_id, :code, :discount)"
                ),
                {
                    "user_id": user_id,
                    "code": "INV_" + req.invitation_code,
                    "discount": 1500,
                },
            )

            # 招待した人にもRewardを付与
            conn.execute(
                text(
                    "INSERT INTO coupons (user_id, code, discount) VALUES (:user_id, CONCAT(:code_prefix, '_', FLOOR(UNIX_TIMESTAMP(NOW(3))*1000)), :discount)"
                ),
                {
                    "user_id": inviter.id,
                    "code_prefix": "RWD_" + req.invitation_code,
                    "discount": 1000,
                },
            )

    response.set_cookie(key="app_session", value=access_token, path="/")
    return AppPostUsersResponse(id=user_id, invitation_code=invitation_code)


class AppPostPaymentMethodsRequest(BaseModel):
    token: str


@router.post("/payment-methods", status_code=HTTPStatus.NO_CONTENT)
def app_post_payment_methods(
    user: Annotated[User, Depends(app_auth_middleware)],
    req: AppPostPaymentMethodsRequest,
) -> None:
    if req.token == "":
        raise HTTPException(
            status_code=HTTPStatus.BAD_REQUEST, detail="token is required but was empty"
        )

    with engine.begin() as conn:
        conn.execute(
            text(
                "INSERT INTO payment_tokens (user_id, token) VALUES (:user_id, :token)"
            ),
            {"user_id": user.id, "token": req.token},
        )


class Coordinate(BaseModel):
    latitude: int
    longitude: int


class GetAppRidesResponseItemChair(BaseModel):
    id: str
    owner: str
    name: str
    model: str


class GetAppRidesResponseItem(BaseModel):
    id: str
    pickup_coordinate: Coordinate
    destination_coordinate: Coordinate
    chair: GetAppRidesResponseItemChair
    fare: int
    evaluation: int
    requested_at: int
    completed_at: int


class GetAppRidesResponse(BaseModel):
    rides: list[GetAppRidesResponseItem]


@router.get("/rides")
def app_get_rides(
    user: Annotated[User, Depends(app_auth_middleware)],
) -> GetAppRidesResponse:
    with engine.begin() as conn:
        rows = conn.execute(
            text(
                "SELECT * FROM rides WHERE user_id = :user_id ORDER BY created_at DESC"
            ),
            {"user_id": user.id},
        ).fetchall()
        rides = [Ride.model_validate(row) for row in rows]

        items = []
        for ride in rides:
            status = get_latest_ride_status(conn, ride.id)
            if status != "COMPLETED":
                continue

            fare = calculate_discounted_fare(
                conn,
                user.id,
                ride,
                ride.pickup_latitude,
                ride.pickup_longitude,
                ride.destination_latitude,
                ride.destination_longitude,
            )

            row = conn.execute(
                text("SELECT * FROM chairs WHERE id = :id"), {"id": ride.chair_id}
            ).fetchone()
            if row is None:
                raise HTTPException(status_code=HTTPStatus.INTERNAL_SERVER_ERROR)
            chair = Chair.model_validate(row)

            row = conn.execute(
                text("SELECT * FROM owners WHERE id = :id"), {"id": chair.owner_id}
            ).fetchone()
            if row is None:
                raise HTTPException(status_code=HTTPStatus.INTERNAL_SERVER_ERROR)
            owner = Owner.model_validate(row)

            item = GetAppRidesResponseItem(
                id=ride.id,
                pickup_coordinate=Coordinate(
                    latitude=ride.pickup_latitude, longitude=ride.pickup_longitude
                ),
                destination_coordinate=Coordinate(
                    latitude=ride.destination_latitude,
                    longitude=ride.destination_longitude,
                ),
                chair=GetAppRidesResponseItemChair(
                    id=chair.id, owner=owner.name, name=chair.name, model=chair.model
                ),
                fare=fare,
                evaluation=ride.evaluation,  # type: ignore[arg-type]
                requested_at=timestamp_millis(ride.created_at),
                completed_at=timestamp_millis(ride.updated_at),
            )
            items.append(item)

    return GetAppRidesResponse(rides=items)


class AppPostRidesRequest(BaseModel):
    pickup_coordinate: Coordinate | None = None
    destination_coordinate: Coordinate | None = None


class AppPostRidesResponse(BaseModel):
    ride_id: str
    fare: int


def get_latest_ride_status(conn: Connection, ride_id: str) -> str:
    status = conn.execute(
        text(
            "SELECT status FROM ride_statuses WHERE ride_id = :ride_id ORDER BY created_at DESC LIMIT 1"
        ),
        {"ride_id": ride_id},
    ).scalar()

    if status is None:
        raise HTTPException(
            status_code=HTTPStatus.INTERNAL_SERVER_ERROR, detail="no rows in result set"
        )

    assert isinstance(status, str)
    return status


@router.post("/rides", status_code=HTTPStatus.ACCEPTED)
def app_post_rides(
    user: Annotated[User, Depends(app_auth_middleware)],
    req: AppPostRidesRequest,
) -> AppPostRidesResponse:
    if req.pickup_coordinate is None or req.destination_coordinate is None:
        raise HTTPException(
            status_code=HTTPStatus.BAD_REQUEST,
            detail="required fields(pickup_coordinate, destination_coordinate) are empty",
        )

    ride_id = str(ULID())
    with engine.begin() as conn:
        rows = conn.execute(
            text("SELECT * FROM rides WHERE user_id = :user_id"), {"user_id": user.id}
        ).fetchall()
        rides = [Ride.model_validate(row) for row in rows]

        continuing_ride_count = 0
        for ride in rides:
            status = get_latest_ride_status(conn, ride.id)
            if status != "COMPLETED":
                continuing_ride_count += 1

        if continuing_ride_count > 0:
            raise HTTPException(
                status_code=HTTPStatus.CONFLICT, detail="ride already exists"
            )

        conn.execute(
            text(
                "INSERT INTO rides (id, user_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude) VALUES (:id, :user_id, :pickup_latitude, :pickup_longitude, :destination_latitude, :destination_longitude)"
            ),
            {
                "id": ride_id,
                "user_id": user.id,
                "pickup_latitude": req.pickup_coordinate.latitude,
                "pickup_longitude": req.pickup_coordinate.longitude,
                "destination_latitude": req.destination_coordinate.latitude,
                "destination_longitude": req.destination_coordinate.longitude,
            },
        )

        conn.execute(
            text(
                "INSERT INTO ride_statuses (id, ride_id, status) VALUES (:id, :ride_id, :status)"
            ),
            {"id": str(ULID()), "ride_id": ride_id, "status": "MATCHING"},
        )

        ride_count = conn.execute(
            text("SELECT COUNT(*) FROM rides WHERE user_id = :user_id"),
            {"user_id": user.id},
        ).scalar()

        if ride_count == 1:
            # 初回利用で、初回利用クーポンがあれば必ず使う
            coupon = conn.execute(
                text(
                    "SELECT * FROM coupons WHERE user_id = :user_id AND code = 'CP_NEW2024' AND used_by IS NULL FOR UPDATE"
                ),
                {"user_id": user.id},
            ).fetchone()

            if coupon:
                conn.execute(
                    text(
                        "UPDATE coupons SET used_by = :ride_id WHERE user_id = :user_id AND code = 'CP_NEW2024'"
                    ),
                    {"ride_id": ride_id, "user_id": user.id},
                )
            else:
                # 無ければ他のクーポンを付与された順番に使う
                coupon = conn.execute(
                    text(
                        "SELECT * FROM coupons WHERE user_id = :user_id AND used_by IS NULL ORDER BY created_at LIMIT 1 FOR UPDATE"
                    ),
                    {"user_id": user.id},
                ).fetchone()
                if coupon:
                    conn.execute(
                        text(
                            "UPDATE coupons SET used_by = :ride_id WHERE user_id = :user_id AND code = :code"
                        ),
                        {"ride_id": ride_id, "user_id": user.id, "code": coupon.code},
                    )
        else:
            # 他のクーポンを付与された順番に使う
            coupon = conn.execute(
                text(
                    "SELECT * FROM coupons WHERE user_id = :user_id AND used_by IS NULL ORDER BY created_at LIMIT 1 FOR UPDATE"
                ),
                {"user_id": user.id},
            ).fetchone()
            if coupon:
                conn.execute(
                    text(
                        "UPDATE coupons SET used_by = :ride_id WHERE user_id = :user_id AND code = :code"
                    ),
                    {"ride_id": ride_id, "user_id": user.id, "code": coupon.code},
                )

        row = conn.execute(
            text("SELECT * FROM rides WHERE id = :ride_id"), {"ride_id": ride_id}
        ).fetchone()
        ride = Ride.model_validate(row)

        fare = calculate_discounted_fare(
            conn,
            user.id,
            ride,
            req.pickup_coordinate.latitude,
            req.pickup_coordinate.longitude,
            req.destination_coordinate.latitude,
            req.destination_coordinate.longitude,
        )

    return AppPostRidesResponse(ride_id=ride_id, fare=fare)


class AppPostRidesEstimatedFareRequest(BaseModel):
    pickup_coordinate: Coordinate | None = None
    destination_coordinate: Coordinate | None = None


class AppPostRidesEstimatedFareResponse(BaseModel):
    fare: int
    discount: int


@router.post(
    "/rides/estimated-fare",
    status_code=HTTPStatus.OK,
)
def app_post_rides_estimated_fare(
    user: Annotated[User, Depends(app_auth_middleware)],
    req: AppPostRidesEstimatedFareRequest,
) -> AppPostRidesEstimatedFareResponse:
    if req.pickup_coordinate is None or req.destination_coordinate is None:
        raise HTTPException(
            status_code=HTTPStatus.BAD_REQUEST,
            detail="required fields(pickup_coordinate, destination_coordinate) are empty",
        )

    with engine.begin() as conn:
        discounted = calculate_discounted_fare(
            conn,
            user.id,
            None,
            req.pickup_coordinate.latitude,
            req.pickup_coordinate.longitude,
            req.destination_coordinate.latitude,
            req.destination_coordinate.longitude,
        )

        return AppPostRidesEstimatedFareResponse(
            fare=discounted,
            discount=calculate_fare(
                req.pickup_coordinate.latitude,
                req.pickup_coordinate.longitude,
                req.destination_coordinate.latitude,
                req.destination_coordinate.longitude,
            )
            - discounted,
        )


class AppPostRideEvaluationRequest(BaseModel):
    evaluation: int


class AppPostRideEvaluationResponse(BaseModel):
    completed_at: int


@router.post(
    "/rides/{ride_id}/evaluation",
    status_code=HTTPStatus.OK,
)
def app_post_ride_evaluation(
    _user: Annotated[User, Depends(app_auth_middleware)],
    req: AppPostRideEvaluationRequest,
    ride_id: str,
) -> AppPostRideEvaluationResponse:
    if req.evaluation < 1 or req.evaluation > 5:
        raise HTTPException(
            status_code=HTTPStatus.BAD_REQUEST,
            detail="evaluation must be between 1 and 5",
        )

    with engine.begin() as conn:
        row = conn.execute(
            text("SELECT * FROM rides WHERE id = :ride_id"), {"ride_id": ride_id}
        ).fetchone()

        if row is None:
            raise HTTPException(
                status_code=HTTPStatus.NOT_FOUND, detail="ride not found"
            )
        ride = Ride.model_validate(row)
        status = get_latest_ride_status(conn, ride.id)

        if status != "ARRIVED":
            raise HTTPException(
                status_code=HTTPStatus.BAD_REQUEST, detail="not arrived yet"
            )

        result = conn.execute(
            text("UPDATE rides SET evaluation = :evaluation WHERE id = :id"),
            {"evaluation": req.evaluation, "id": ride_id},
        )
        if result.rowcount == 0:
            raise HTTPException(
                status_code=HTTPStatus.NOT_FOUND, detail="ride not found"
            )

        conn.execute(
            text(
                "INSERT INTO ride_statuses (id, ride_id, status) VALUES (:id, :ride_id, :status)"
            ),
            {"id": str(ULID()), "ride_id": ride_id, "status": "COMPLETED"},
        )

        row = conn.execute(
            text("SELECT * FROM rides WHERE id = :id"), {"id": ride_id}
        ).fetchone()
        if row is None:
            raise HTTPException(
                status_code=HTTPStatus.NOT_FOUND, detail="ride not found"
            )
        ride = Ride.model_validate(row)

        row = conn.execute(
            text("SELECT * FROM payment_tokens WHERE user_id = :user_id"),
            {"user_id": ride.user_id},
        ).fetchone()
        if row is None:
            raise HTTPException(
                status_code=HTTPStatus.BAD_REQUEST,
                detail="payment token not registered",
            )
        payment_token = PaymentToken.model_validate(row)

        fare = calculate_discounted_fare(
            conn,
            ride.user_id,
            ride,
            ride.pickup_latitude,
            ride.pickup_longitude,
            ride.destination_latitude,
            ride.destination_longitude,
        )

        payment_gateway_request = PaymentGatewayPostPaymentRequest(amount=fare)
        payment_gateway_url = conn.execute(
            text("SELECT value FROM settings WHERE name = 'payment_gateway_url'")
        ).scalar()
        if not isinstance(payment_gateway_url, str):
            raise HTTPException(status_code=HTTPStatus.INTERNAL_SERVER_ERROR)

        def retrieve_rides_order_by_created_at_asc() -> list[Ride]:
            rows = conn.execute(
                text(
                    "SELECT * FROM rides WHERE user_id = :user_id ORDER BY created_at ASC",
                ),
                {"user_id": ride.user_id},
            ).fetchall()
            return [Ride.model_validate(r) for r in rows]

        try:
            request_payment_gateway_post_payment(
                payment_gateway_url,
                payment_token.token,
                payment_gateway_request,
                retrieve_rides_order_by_created_at_asc,
            )
        except UpstreamError as e:
            raise HTTPException(status_code=HTTPStatus.BAD_GATEWAY, detail=str(e))

        response = AppPostRideEvaluationResponse(
            completed_at=timestamp_millis(ride.updated_at)
        )
    return response


class AppGetNotificationResponseChairStats(BaseModel):
    total_rides_count: int
    total_evaluation_avg: float


class AppGetNotificationResponseChair(BaseModel):
    id: str
    name: str
    model: str
    stats: AppGetNotificationResponseChairStats


class AppGetNotificationResponseData(BaseModel):
    ride_id: str
    pickup_coordinate: Coordinate
    destination_coordinate: Coordinate
    fare: int
    status: str
    chair: AppGetNotificationResponseChair | None = None
    created_at: int
    updated_at: int


class AppGetNotificationResponse(BaseModel):
    data: AppGetNotificationResponseData | None = None
    retry_after_ms: int | None = None


@router.get(
    "/notification",
    status_code=HTTPStatus.OK,
    response_model_exclude_none=True,
)
def app_get_notification(
    user: Annotated[User, Depends(app_auth_middleware)],
) -> AppGetNotificationResponse:
    with engine.begin() as conn:
        row = conn.execute(
            text(
                "SELECT * FROM rides WHERE user_id = :user_id ORDER BY created_at DESC LIMIT 1"
            ),
            {"user_id": user.id},
        ).fetchone()
        if row is None:
            notification_response = AppGetNotificationResponse(retry_after_ms=30)
            return notification_response

        ride: Ride = Ride.model_validate(row)

        row = conn.execute(
            text(
                "SELECT * FROM ride_statuses WHERE ride_id = :ride_id AND app_sent_at IS NULL ORDER BY created_at ASC LIMIT 1"
            ),
            {"ride_id": ride.id},
        ).fetchone()
        yet_sent_ride_status: RideStatus | None = None
        if row is None:
            status = get_latest_ride_status(conn, ride.id)
        else:
            yet_sent_ride_status = RideStatus.model_validate(row)
            status = yet_sent_ride_status.status

        fare = calculate_discounted_fare(
            conn,
            user.id,
            ride,
            ride.pickup_latitude,
            ride.pickup_longitude,
            ride.destination_latitude,
            ride.destination_longitude,
        )

        notification_response = AppGetNotificationResponse(
            data=AppGetNotificationResponseData(
                ride_id=ride.id,
                pickup_coordinate=Coordinate(
                    latitude=ride.pickup_latitude, longitude=ride.pickup_longitude
                ),
                destination_coordinate=Coordinate(
                    latitude=ride.destination_latitude,
                    longitude=ride.destination_longitude,
                ),
                fare=fare,
                status=status,
                chair=None,
                created_at=timestamp_millis(ride.created_at),
                updated_at=timestamp_millis(ride.updated_at),
            ),
            retry_after_ms=30,
        )

        if ride.chair_id:
            row = conn.execute(
                text("SELECT * FROM chairs WHERE id = :chair_id"),
                {"chair_id": ride.chair_id},
            ).fetchone()
            if row is None:
                raise HTTPException(status_code=HTTPStatus.INTERNAL_SERVER_ERROR)

            chair: Chair = Chair.model_validate(row)

            stats = get_chair_stats(conn, ride.chair_id)

            notification_response.data.chair = AppGetNotificationResponseChair(  # type: ignore[union-attr]
                id=chair.id, name=chair.name, model=chair.model, stats=stats
            )

        if yet_sent_ride_status:
            conn.execute(
                text(
                    "UPDATE ride_statuses SET app_sent_at = CURRENT_TIMESTAMP(6) WHERE id = :yet_send_ride_status_id"
                ),
                {"yet_send_ride_status_id": yet_sent_ride_status.id},
            )

    return notification_response


class RecentRide(BaseModel):
    id: str
    pickup_coordinate: Coordinate
    destination_coordinate: Coordinate
    distance: int
    duration: int
    evaluation: int


def get_chair_stats(
    conn: Connection, chair_id: str
) -> AppGetNotificationResponseChairStats:
    rides = conn.execute(
        text("SELECT * FROM rides WHERE chair_id = :chair_id ORDER BY updated_at DESC"),
        {"chair_id": chair_id},
    ).fetchall()
    total_ride_count = 0
    total_evaluation = 0.0

    for ride in rides:
        rows = conn.execute(
            text(
                "SELECT * FROM ride_statuses WHERE ride_id = :ride_id ORDER BY created_at"
            ),
            {"ride_id": ride.id},
        )
        ride_statuses = [RideStatus.model_validate(row) for row in rows]

        arrived_at = None
        pickuped_at = None
        is_completed = None
        for status in ride_statuses:
            if status.status == "ARRIVED":
                arrived_at = status.created_at
            elif status.status == "CARRYING":
                pickuped_at = status.created_at
            if status.status == "COMPLETED":
                is_completed = True

        if (arrived_at is None) or (pickuped_at is None):
            continue
        if not is_completed:
            continue

        total_ride_count += 1
        total_evaluation += float(ride.evaluation)

    if total_ride_count > 0:
        total_evaluation_avg = total_evaluation / total_ride_count
    else:
        total_evaluation_avg = 0.0

    return AppGetNotificationResponseChairStats(
        total_rides_count=total_ride_count, total_evaluation_avg=total_evaluation_avg
    )


class AppGetNearbyChairsResponseChair(BaseModel):
    id: str
    name: str
    model: str
    current_coordinate: Coordinate


class AppGetNearByChairsResponse(BaseModel):
    chairs: list[AppGetNearbyChairsResponseChair]
    retrieved_at: int


@router.get(
    "/nearby-chairs",
    status_code=HTTPStatus.OK,
)
def app_get_nearby_chairs(
    _user: Annotated[User, Depends(app_auth_middleware)],
    latitude: int,
    longitude: int,
    distance: int = 50,
) -> AppGetNearByChairsResponse:
    coordinate = Coordinate(latitude=latitude, longitude=longitude)
    with engine.begin() as conn:
        chairs = conn.execute(
            text("SELECT * FROM chairs"),
        ).fetchall()

        near_by_chairs = []
        for chair in chairs:
            if not chair.is_active:
                continue

            rows = conn.execute(
                text(
                    "SELECT * FROM rides WHERE chair_id = :chair_id ORDER BY created_at DESC"
                ),
                {"chair_id": chair.id},
            ).fetchall()
            rides = [Ride.model_validate(row) for row in rows]

            skip = False

            for ride in rides:
                # 過去にライドが存在し、かつ、それが完了していない場合はスキップ
                status = get_latest_ride_status(conn, ride.id)
                if status != "COMPLETED":
                    skip = True
                    break

            if skip:
                continue

            # 最新の位置情報を取得
            row = conn.execute(
                text(
                    "SELECT * FROM chair_locations WHERE chair_id = :chair_id ORDER BY created_at DESC LIMIT 1"
                ),
                {"chair_id": chair.id},
            ).fetchone()
            if row is None:
                continue

            chair_location = ChairLocation.model_validate(row)

            if (
                calculate_distance(
                    coordinate.latitude,
                    coordinate.longitude,
                    chair_location.latitude,
                    chair_location.longitude,
                )
                <= distance
            ):
                near_by_chairs.append(
                    AppGetNearbyChairsResponseChair(
                        id=chair.id,
                        name=chair.name,
                        model=chair.model,
                        current_coordinate=Coordinate(
                            latitude=chair_location.latitude,
                            longitude=chair_location.longitude,
                        ),
                    )
                )
        retrieved_at = conn.execute(text("SELECT CURRENT_TIMESTAMP(6)")).scalar()
        assert retrieved_at is not None

    return AppGetNearByChairsResponse(
        chairs=near_by_chairs,
        retrieved_at=timestamp_millis(retrieved_at),
    )


def calculate_discounted_fare(
    conn: Connection,
    user_id: str,
    ride: Ride | None,
    pickup_latitude: int,
    pickup_longitude: int,
    dest_latitude: int,
    dest_longitude: int,
) -> int:
    discount = 0

    if ride:
        dest_latitude = ride.destination_latitude
        dest_longitude = ride.destination_longitude
        pickup_latitude = ride.pickup_latitude
        pickup_longitude = ride.pickup_longitude

        # すでにクーポンが紐づいているならそれの割引額を参照
        coupon = conn.execute(
            text("SELECT * FROM coupons WHERE used_by = :ride_id"), {"ride_id": ride.id}
        ).fetchone()
        if coupon:
            discount = coupon.discount
    else:
        # 初回利用クーポンを最優先で使う
        coupon = conn.execute(
            text(
                "SELECT * FROM coupons WHERE user_id = :user_id AND code = 'CP_NEW2024' AND used_by IS NULL"
            ),
            {"user_id": user_id},
        ).fetchone()

        if coupon is None:
            # 無いなら他のクーポンを付与された順番に使う
            coupon = conn.execute(
                text(
                    "SELECT * FROM coupons WHERE user_id = :user_id AND used_by IS NULL ORDER BY created_at LIMIT 1"
                ),
                {"user_id": user_id},
            ).fetchone()

        if coupon:
            discount = coupon.discount

    metered_fare = FARE_PER_DISTANCE * calculate_distance(
        dest_latitude, dest_longitude, pickup_latitude, pickup_longitude
    )

    discounted_metered_fare = max(metered_fare - discount, 0)
    return INITIAL_FARE + discounted_metered_fare
