"""
以下の移植
https://github.com/isucon/isucon14/blob/main/webapp/go/app_handlers.go

TODO: このdocstringを消す
"""

from http import HTTPStatus

from fastapi import APIRouter, Depends, HTTPException, Response
from pydantic import BaseModel
from sqlalchemy import text
from ulid import ULID

from ..middlewares import app_auth_middleware
from ..models import Chair, Ride, User
from ..sql import engine
from ..utils import secure_random_str
from .owners import calculate_sale, fare_per_distance, initial_fare

router = APIRouter(prefix="/api/app")


@router.post("/register", status_code=201)
def app_post_register():
    user_id = str(ULID())
    return {"id": user_id}


class AppPostPaymentMethodsRequest(BaseModel):
    token: str


@router.post("/payment-methods", status_code=HTTPStatus.NO_CONTENT)
def app_post_payment_methods(
    req: AppPostPaymentMethodsRequest, user: User = Depends(app_auth_middleware)
):
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


@router.post("/requests", status_code=202)
def app_post_requests():
    request_id = str(ULID())
    return {"request_id": request_id}


@router.get("/users", status_code=200)
def app_get_users():
    pass


class Coordinate(BaseModel):
    latitude: int
    longitude: int


class AppPostRidesRequest(BaseModel):
    pickup_coordinate: Coordinate | None = None
    destination_coordinate: Coordinate | None = None


class AppPostRidesResponse(BaseModel):
    ride_id: str
    fare: int


def get_latest_ride_status(conn, ride_id: str) -> str:
    row = conn.execute(
        text(
            "SELECT status FROM ride_statuses WHERE ride_id = :ride_id ORDER BY created_at DESC LIMIT 1"
        ),
        {"ride_id": ride_id},
    ).fetchone()

    if not row:
        return ""

    return row.status


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
def app_get_rides(user: User = Depends(app_auth_middleware)):
    with engine.begin() as conn:
        rows = conn.execute(
            text(
                "SELECT * FROM rides WHERE user_id = :user_id ORDER BY created_at DESC"
            ),
            {"user_id": user.id},
        ).fetchall()
        rides = [Ride(**row._mapping) for row in rows]

    items = []
    for ride in rides:
        status = get_latest_ride_status(conn, ride.id)
        if status != "COMPLETED":
            continue

        with engine.begin() as conn:
            row = conn.execute(
                text("SELECT * FROM chairs WHERE id = :id"), {"id": ride.chair_id}
            ).fetchone()
            if row is None:
                raise HTTPException(status_code=HTTPStatus.INTERNAL_SERVER_ERROR)
            chair = Chair(**row._mapping)

        with engine.begin() as conn:
            row = conn.execute(
                text("SELECT * FROM owners WHERE id = :id"), {"id": chair.owner_id}
            ).fetchone()
            if row is None:
                raise HTTPException(status_code=HTTPStatus.INTERNAL_SERVER_ERROR)
            owner = Chair(**row._mapping)

        # TODO: 参照実装みたいにpartialに作るべき？
        item = GetAppRidesResponseItem(
            id=ride.id,
            pickup_coordinate=Coordinate(
                latitude=ride.pickup_latitude, longitude=ride.pickup_longitude
            ),
            destination_coordinate=Coordinate(
                latitude=ride.destination_latitude, longitude=ride.destination_longitude
            ),
            chair=GetAppRidesResponseItemChair(
                id=chair.id, owner=owner.name, name=chair.name, model=chair.model
            ),
            fare=calculate_sale(ride),
            # TODO: 型エラーを修正
            evaluation=ride.evaluation,  # type: ignore[arg-type]
            requested_at=int(ride.created_at.timestamp() * 1000),
            completed_at=int(ride.updated_at.timestamp() * 1000),
        )
        items.append(item)

    return GetAppRidesResponse(rides=items)


@router.post("/rides", status_code=202)
def app_post_rides(
    r: AppPostRidesRequest, user: User = Depends(app_auth_middleware)
) -> AppPostRidesResponse:
    if r.pickup_coordinate is None or r.destination_coordinate is None:
        raise HTTPException(
            status_code=400,
            detail="required fields(pickup_coordinate, destination_coordinate) are empty",
        )

    ride_id = str(ULID())
    with engine.begin() as conn:
        rides = conn.execute(
            text("SELECT * FROM rides WHERE user_id = :user_id"), {"user_id": user.id}
        ).fetchall()

        continuing_ride_count: int = 0
        for ride in rides:
            status = get_latest_ride_status(conn, ride.id)
            if status != "COMPLETED" and status != "CANCELED":
                continuing_ride_count += 1

        if continuing_ride_count > 0:
            raise HTTPException(status_code=409, detail="ride already exists")

        conn.execute(
            text(
                "INSERT INTO rides (id, user_id, pickup_latitude, pickup_longitude, destination_latitude, destination_longitude) VALUES (:id, :user_id, :pickup_latitude, :pickup_longitude, :destination_latitude, :destination_longitude)"
            ),
            {
                "id": ride_id,
                "user_id": user.id,
                "pickup_latitude": r.pickup_coordinate.latitude,
                "pickup_longitude": r.pickup_coordinate.longitude,
                "destination_latitude": r.destination_coordinate.latitude,
                "destination_longitude": r.destination_coordinate.longitude,
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
                coupon = conn.execute(  # type: ignore
                    text(
                        "SELECT * FROM coupons WHERE user_id = :user_id AND used_by IS NULL ORDER BY created_at LIMIT 1 FOR UPDATE"
                    ),
                    {"user_id": user.id},
                )
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
        ride: Ride = Ride(**row._mapping)  # type: ignore

        fare = calculate_discounted_fare(
            conn,
            user.id,
            ride,  # type: ignore
            r.pickup_coordinate.latitude,
            r.pickup_coordinate.longitude,
            r.destination_coordinate.latitude,
            r.destination_coordinate.longitude,
        )

    return AppPostRidesResponse(ride_id=ride_id, fare=fare)


class AppPostRidesEstimatedFareRequest(BaseModel):
    pickup_coordinate: Coordinate | None = None
    destination_coordinate: Coordinate | None = None


class AppPostRidesEstimatedFareResponse(BaseModel):
    fare: int
    discount: int


@router.post(
    "/api/app/rides/estimated-fare",
    response_model=AppPostRidesEstimatedFareResponse,
    status_code=200,
)
def app_post_rides_estimated_fare(
    r: AppPostRidesEstimatedFareRequest, user: User = Depends(app_auth_middleware)
) -> AppPostRidesEstimatedFareResponse:
    if r.pickup_coordinate is None or r.destination_coordinate is None:
        raise HTTPException(
            status_code=400,
            detail="required fields(pickup_coordinate, destination_coordinate) are empty",
        )

    with engine.begin() as conn:
        discounted = calculate_discounted_fare(
            conn,
            user.id,
            None,
            r.pickup_coordinate.latitude,
            r.pickup_coordinate.longitude,
            r.destination_coordinate.latitude,
            r.destination_coordinate.longitude,
        )

    return AppPostRidesEstimatedFareResponse(
        fare=discounted,
        discount=calculate_discounted_fare(
            conn,
            user.id,
            None,
            r.pickup_coordinate.latitude,
            r.pickup_coordinate.longitude,
            r.destination_coordinate.latitude,
            r.destination_coordinate.longitude,
        )
        - discounted,
    )


def calculate_distance(
    a_latitude: int, a_longitude: int, b_latitude: int, b_longitude: int
) -> int:
    return abs(a_latitude - b_latitude) + abs(a_longitude - b_longitude)


class AppPostRideEvaluationRequest(BaseModel):
    evaluation: int


class AppPostRideEvaluationResponse(BaseModel):
    completed_at: int


@router.post(
    "/rides/{ride_id}/evaluation",
    response_model=AppPostRideEvaluationResponse,
    status_code=200,
)
def app_post_ride_evaluatation(
    r: AppPostRideEvaluationRequest, ride_id: str
) -> AppPostRideEvaluationResponse:
    if r.evaluation < 1 or r.evaluation > 5:
        raise HTTPException(
            status_code=400, detail="evaluation must be between 1 and 5"
        )

    with engine.begin() as conn:
        row = conn.execute(
            text("SELECT * FROM rides WHERE id = :ride_id"), {"ride_id": ride_id}
        ).fetchone()

        if not row:
            raise HTTPException(status_code=404, detail="ride not found")
        ride = Ride(**row._mapping)
        status = get_latest_ride_status(conn, ride.id)

        if status != "ARRIVED":
            raise HTTPException(status_code=400, detail="not arrived yet")

        # TODO: write a rest here
        response = AppPostRideEvaluationResponse(completed_at=10000)
    return response


class AppPostUsersRequest(BaseModel):
    username: str
    firstname: str
    lastname: str
    date_of_birth: str
    invitation_code: str | None = None


class AppPostUsersResponse(BaseModel):
    id: str
    invitation_code: str


@router.post("/users", response_model=AppPostUsersResponse, status_code=201)
def app_post_users(r: AppPostUsersRequest, response: Response) -> AppPostUsersResponse:
    user_id = str(ULID())
    access_token = secure_random_str(32)
    invitation_code = secure_random_str(15)

    # start transaction
    with engine.begin() as conn:
        conn.execute(
            text(
                "INSERT INTO users (id, username, firstname, lastname, date_of_birth, access_token, invitation_code) VALUES (:id, :username, :firstname, :lastname, :date_of_birth, :access_token, :invitation_code)"
            ),
            {
                "id": user_id,
                "username": r.username,
                "firstname": r.firstname,
                "lastname": r.lastname,
                "date_of_birth": r.date_of_birth,
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
        if r.invitation_code:
            # 招待する側の招待数をチェック
            coupons = conn.execute(
                text("SELECT * FROM coupons WHERE code = :code FOR UPDATE"),
                {"code": "INV_" + r.invitation_code},
            ).fetchall()

            if len(coupons) >= 3:
                raise HTTPException(
                    status_code=400, detail="この招待コードは使用できません"
                )

            # ユーザーチェック
            inviter = conn.execute(
                text("SELECT * FROM users WHERE invitation_code = :invitation_code"),
                {"invitation_code": r.invitation_code},
            ).fetchone()

            if not inviter:
                raise HTTPException(
                    status_code=400, detail="この招待コードは使用できません。"
                )

            # 招待クーポン付与
            conn.execute(
                text(
                    "INSERT INTO coupons (user_id, code, discount) VALUES (:user_id, :code, :discount)"
                ),
                {
                    "user_id": user_id,
                    "code": "INV_" + r.invitation_code,
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
                    "code": "RWD_" + r.invitation_code,
                    "discount": 1000,
                },
            )

    response.set_cookie(key="app_session", value=access_token, path="/")
    return AppPostUsersResponse(id=user_id, invitation_code=invitation_code)


def calculate_discounted_fare(
    conn,
    user_id: str,
    ride: Ride | None,
    pickup_latitude: int,
    pickup_longitude: int,
    dest_latitude: int,
    dest_longitude: int,
) -> int:
    discount: int = 0

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

        if not coupon:
            # 無いなら他のクーポンを付与された順番に使う
            coupon = conn.execute(
                text(
                    "SELECT * FROM coupons WHERE user_id = :user_id AND used_by IS NULL ORDER BY created_at LIMIT 1"
                ),
                {"user_id": user_id},
            ).fetchone()

        if coupon:
            discount = coupon.discount

    metered_fare: int = fare_per_distance * calculate_distance(
        dest_latitude, dest_longitude, pickup_latitude, pickup_longitude
    )

    discounted_metered_fare: int = max(metered_fare - discount, 0)
    return initial_fare + discounted_metered_fare


class RecentRide(BaseModel):
    id: str
    pickup_coordinate: Coordinate
    destination_coordinate: Coordinate
    distance: int
    duration: int
    evaluation: int


class AppChairStats(BaseModel):
    # 最近の乗車履歴
    recent_rides: list[RecentRide]

    # 累計の情報
    total_rides_count: int
    total_evaluation_avg: float


class AppChair(BaseModel):
    id: str
    name: str
    model: str
    stats: AppChairStats


class AppGetRideResponse(BaseModel):
    id: str
    pickup_coordinate: Coordinate
    destination_coordinate: Coordinate
    status: str
    chair: AppChair | None = None
    created_at: int
    updated_at: int


@router.get(
    "/rides/{ride_id}",
    response_model=AppGetRideResponse,
    status_code=200,
    response_model_exclude_none=True,
)
def app_get_ride(
    ride_id: str, user: User = Depends(app_auth_middleware)
) -> AppGetRideResponse:
    with engine.begin() as conn:
        ride = conn.execute(
            text("SELECT * FROM rides WHERE id = :ride_id"), {"ride_id": ride_id}
        ).fetchone()
        if not ride:
            raise HTTPException(status_code=404, detail="ride not found")
        status = get_latest_ride_status(conn, ride.id)

        response = AppGetRideResponse(
            id=ride.id,
            pickup_coordinate=Coordinate(
                latitude=ride.pickup_latitude, longitude=ride.pickup_longitude
            ),
            destination_coordinate=Coordinate(
                latitude=ride.destination_latitude, longitude=ride.destination_longitude
            ),
            status=status,
            created_at=int(ride.created_at.timestamp() * 1000),
            updated_at=int(ride.updated_at.timestamp() * 1000),
        )

        if ride.chair_id:
            chair = conn.execute(
                text("SELECT * FROM chairs WHERE id = :chair_id"),
                {"chair_id": ride.chair_id},
            ).fetchone()

            # TODO: stats = get_chair_stats(chair.id)
            stats = AppChairStats(
                recent_rides=[], total_rides_count=1, total_evaluation_avg=0.1
            )
            response.chair = AppChair(
                id=chair.id,  # type: ignore
                name=chair.name,  # type: ignore
                model=chair.model,  # type: ignore
                stats=stats,  # type: ignore
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


class AppGetNotificationResponse(BaseModel):
    ride_id: str
    pickup_coordinate: Coordinate
    destination_coordinate: Coordinate
    fare: int
    status: str
    chair: AppGetNotificationResponseChair | None = None
    created_at: int
    updated_at: int


@router.get(
    "/notification",
    response_model=AppGetNotificationResponse,
    status_code=200,
    response_model_exclude_none=True,
)
def app_get_notification(
    response: Response, user: User = Depends(app_auth_middleware)
) -> AppGetNotificationResponse | Response:
    with engine.begin() as conn:
        row = conn.execute(
            text(
                "SELECT * FROM rides WHERE user_id = :user_id ORDER BY created_at DESC LIMIT 1"
            ),
            {"user_id": user.id},
        ).fetchone()
        if not row:
            response.status_code = HTTPStatus.NO_CONTENT
            return response

        ride: Ride = Ride(**row._mapping)
        status = get_latest_ride_status(conn, ride.id)

        notification_response = AppGetNotificationResponse(
            ride_id=ride.id,
            pickup_coordinate=Coordinate(latitude=0, longitude=0),
            destination_coordinate=Coordinate(latitude=10, longitude=10),
            fare=100,
            status=status,
            chair=None,
            created_at=1000,
            updated_at=10000,
        )

        if ride.chair_id:
            row = conn.execute(
                text("SELECT * FROM chairs WHERE id = :chair_id"),
                {"chair_id": ride.chair_id},
            ).fetchone()

            chair: Chair = Chair(**row._mapping)

            stats = get_chair_stats(conn, ride.chair_id)

            notification_response.chair = AppGetNotificationResponseChair(
                id=chair.id, name=chair.name, model=chair.model, stats=chair.stats
            )

        # TODO: check the chair is here
    return notification_response


def get_chair_stats(conn, chair_id: str) -> AppGetNotificationResponseChairStats:
    rides = engine.execute(
        text("SELECT * FROM rides WHERE chair_id = :chair_id ORDER BY updated_at DESC"),
        {"chair_id": chair_id},
    ).fetchall()
    total_ride_count: int = len(rides)
    total_evaluation: float = 0.0

    for ride in rides:
        chair_locations = conn.execute(
            text(
                "SELECT * FROM chair_locations WHERE chair_id = :chair_id AND created_at BETWEEN :created_at AND :updated_at ORDER BY created_at"
            ),
            {
                "chair_id": chair_id,
                "created_at": ride.created_at,
                "updated_at": ride.updated_at,
            },
        )


class AppGetNearByChairsResponse(BaseModel):
    chairs: list[AppChair]
    retrieved_at: int


class AppGetNearbyChairsResponseChair(BaseModel):
    id: str
    name: str
    model: str
    current_coordinate: Coordinate


@router.get(
    "/nearby-chairs",
    response_model=AppGetNearByChairsResponse,
    status_code=200,
)
def app_get_near_by_chairs(latitude: int, longitude: int, distance: int = 50):
    coordinate = Coordinate(latitude, longitude)
    with engine.begin() as conn:
        chairs = conn.execute(
            text("SELECT * FROM chairs"),
        ).fetchall()

        near_by_chairs = []
        for chair in chairs:
            ride = conn.execute(
                "SELECT * FROM rides WHERE chair_id = :chair_id ORDER BY created_at DESC LIMIT 1",
                {"chair_id": chair.id},
            ).fetchone()

            if ride:
                # 過去にライドが存在し、かつ、それが完了していない場合はスキップ
                status = get_latest_ride_status(conn, ride.id)
                if status != "COMPLETED":
                    continue

            # 5分以内に更新されている最新の位置情報を取得
            chair_location = conn.execute(
                "SELECT * FROM chair_locations WHERE chair_id = :chair_id AND created_at > DATE_SUB(CURRENT_TIMESTAMP(6), INTERVAL 5 MINUTE) ORDER BY created_at DESC LIMIT 1",
                {"chair_id": chair.id},
            ).fetchone()

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

    return AppGetNearByChairsResponse(chairs=near_by_chairs, retrieved_at=retrieved_at)
