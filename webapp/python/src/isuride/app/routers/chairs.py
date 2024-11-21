"""
以下の移植
https://github.com/isucon/isucon14/blob/main/webapp/go/chair_handlers.go

TODO: このdocstringを消す
"""

from fastapi import APIRouter, Depends, HTTPException, Response, status
from pydantic import BaseModel
from sqlalchemy import text
from ulid import ULID

from ..middlewares import chair_auth_middleware
from ..models import Chair, ChairLocation, Owner, Ride, User
from ..sql import engine
from ..utils import secure_random_str
from .apps import get_latest_ride_status

router = APIRouter(prefix="/api/chair")


class ChairPostChairsRequest(BaseModel):
    name: str
    model: str
    chair_register_token: str


class ChairPostChairsResponse(BaseModel):
    id: str
    owner_id: str


@router.post("/chairs", status_code=status.HTTP_201_CREATED)
def chair_post_chairs(
    req: ChairPostChairsRequest, resp: Response
) -> ChairPostChairsResponse:
    if req.name == "" or req.model == "" or req.chair_register_token == "":
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="some of required fields(name, model, chair_register_token) are empty",
        )

    with engine.begin() as conn:
        row = conn.execute(
            text(
                "SELECT * FROM owners WHERE chair_register_token = :chair_register_token"
            ),
            {"chair_register_token": req.chair_register_token},
        ).fetchone()
        if row is None:
            raise HTTPException(
                status_code=status.UNAUTHORIZED, detail="invalid chair_register_token"
            )
        owner = Owner(**row._mapping)

    chair_id = str(ULID())
    access_token = secure_random_str(32)

    with engine.begin() as conn:
        conn.execute(
            text(
                "INSERT INTO chairs (id, owner_id, name, model, is_active, access_token) VALUES (:id, :owner_id, :name, :model, :is_active, :access_token)",
            ),
            {
                "id": chair_id,
                "owner_id": owner.id,
                "name": req.name,
                "model": req.model,
                "is_active": False,
                "access_token": access_token,
            },
        )

    resp.set_cookie(path="/", key="chair_session", value=access_token)
    return ChairPostChairsResponse(id=chair_id, owner_id=owner.id)


class PostChairActivityRequest(BaseModel):
    is_active: bool


@router.post("/activity", status_code=status.HTTP_204_NO_CONTENT)
def chair_post_activity(
    req: PostChairActivityRequest, chair: Chair = Depends(chair_auth_middleware)
):
    with engine.begin() as conn:
        conn.execute(
            text("UPDATE chairs SET is_active = :is_active WHERE id = :id"),
            {"is_active": req.is_active, "id": chair.id},
        )


# TODO: Requestの構造体がないの、紛らわしいので要検討
class Coordinate(BaseModel):
    latitude: int
    longitude: int


class ChairPostCoordinateResponse(BaseModel):
    recorded_at: int


@router.post("/coordinate")
def chair_post_coordinate(
    req: Coordinate, chair: Chair = Depends(chair_auth_middleware)
):
    with engine.begin() as conn:
        chair_location_id = str(ULID())
        conn.execute(
            text(
                "INSERT INTO chair_locations (id, chair_id, latitude, longitude) VALUES (:id, :chair_id, :latitude, :longitude)"
            ),
            {
                "id": chair_location_id,
                "chair_id": chair.id,
                "latitude": req.latitude,
                "longitude": req.longitude,
            },
        )

        row = conn.execute(
            text("SELECT * FROM chair_locations WHERE id = :id"),
            {"id": chair_location_id},
        ).fetchone()
        if row is None:
            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR)
        location = ChairLocation(**row._mapping)

        row = conn.execute(
            text(
                "SELECT * FROM rides WHERE chair_id = :chair_id ORDER BY updated_at DESC LIMIT 1"
            ),
            {"chair_id": chair.id},
        ).fetchone()
        if row is None:
            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR)

        ride = Ride(**row)
        ride_status = get_latest_ride_status(conn, ride_id=ride.id)
        if ride_status != "COMPLETED" and ride_status != "CANCELLED":
            if (
                req.latitude == ride.pickup_latitude
                and req.longitude == ride.pickup_longitude
                and ride_status == "ENROUTE"
            ):
                conn.execute(
                    text(
                        "INSERT INTO ride_statuses (id, ride_id, status) VALUES (:id, :ride_id, :status)"
                    ),
                    {"id": str(ULID()), "ride_id": ride.id, "status": "PICKUP"},
                )

            if (
                req.latitude == ride.destination_latitude
                and req.longitude == ride.destination_longitude
                and ride_status == "CARRYING"
            ):
                conn.execute(
                    text(
                        "INSERT INTO ride_statuses (id, ride_id, status) VALUES (:id, :ride_id, :status) "
                    ),
                    {"id": str(ULID()), "ride_id": ride.id, "status": "ARRIVED"},
                )

    return ChairPostCoordinateResponse(
        recorded_at=int(location.created_at.timestamp() * 1000)
    )


class SimpleUser(BaseModel):
    id: str
    name: str


class ChairGetNotificationResponse(BaseModel):
    ride_id: str
    user: SimpleUser
    pickup_coordinate: Coordinate
    destination_coordinate: Coordinate
    status: str


@router.get("/notification", status_code=204)
def chair_get_notification(chair: Chair = Depends(chair_auth_middleware)):
    with engine.begin() as conn:
        conn.execute(
            text("SELECT * FROM chairs WHERE id = :id FOR UPDATE"), {"id": chair.id}
        )

        found = True
        ride_status = ""
        row = conn.execute(
            text(
                "SELECT * FROM rides WHERE chair_id = :chair_id ORDER BY updated_at DESC LIMIT 1"
            ),
            {"chair_id": chair.id},
        ).fetchone()
        if row is None:
            found = False

        if found:
            assert row is not None
            ride = Ride(**row)
            ride_status = get_latest_ride_status(conn, ride.id)

        if (not found) or ride_status == "COMPLETED" or ride_status == "CANCELLED":
            # MEMO: 一旦最も待たせているリクエストにマッチさせる実装とする。おそらくもっといい方法があるはず…
            row = conn.execute(
                text(
                    "SELECT * FROM rides WHERE chair_id IS NULL ORDER BY created_at DESC LIMIT 1 FOR UPDATE"
                )
            ).fetchone()
            if row is None:
                raise HTTPException(status_code=status.HTTP_204_NO_CONTENT)
            matched = Ride(**row)

            conn.execute(
                text("UPDATE rides SET chair_id = :chair_id WHERE id = :id"),
                {"chair_id": chair.id, "id": matched.id},
            )

            if not found:
                ride = matched
                ride_status = "MATCHING"

        row = conn.execute(
            text("SELECT * FROM users WHERE id = :id FOR SHARE"), {"id": ride.user_id}
        ).fetchone()
        if row is None:
            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR)
        user = User(**row)

    return ChairGetNotificationResponse(
        ride_id=ride.id,
        user=SimpleUser(id=user.id, name=f"{user.firstname} {user.lastname}"),
        pickup_coordinate=Coordinate(
            latitude=ride.pickup_latitude, longitude=ride.pickup_longitude
        ),
        destination_coordinate=Coordinate(
            latitude=ride.destination_latitude, longitude=ride.destination_longitude
        ),
        status=ride_status,
    )


class PostChairRidesRideIDStatusRequest(BaseModel):
    status: str


@router.post("/rides/{ride_id}/status", status_code=status.HTTP_204_NO_CONTENT)
def chair_post_ride_status(
    ride_id: str,
    req: PostChairRidesRideIDStatusRequest,
    chair: Chair = Depends(chair_auth_middleware),
):
    with engine.begin() as conn:
        row = conn.execute(
            text("SELECT * FROM rides WHERE id = :id FOR UPDATE"), {"id": ride_id}
        ).fetchone()
        if row is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND, detail="ride not found"
            )
        ride = Ride(**row)

        if ride.chair_id != chair.id:
            raise HTTPException(
                status_code=status.BAD_REQUEST, detail="not assigned to this ride"
            )

        match req.status:
            # Deny matching
            case "MATCHING":
                conn.execute(
                    text(
                        "INSERT INTO ride_statuses (id, ride_id, status) VALUES (:id, :ride_id, :status)"
                    ),
                    {"id": str(ULID()), "ride_id": ride.id, "status": "MATCHING"},
                )
            # Accept matching
            case "ENROUTE":
                conn.execute(
                    text(
                        "INSERT INTO ride_statuses (id, ride_id, status) VALUES (:id, :ride_id, :status)"
                    ),
                    {"id": str(ULID()), "ride_id": ride.id, "status": "ENROUTE"},
                )
            # After Picking up user
            case "CARRYING":
                ride_status = get_latest_ride_status(conn, ride.id)
                if ride_status != "PICKUP":
                    raise HTTPException(
                        status_code=status.HTTP_400_BAD_REQUEST,
                        detail="chair has not arrived yet",
                    )
                conn.execute(
                    text(
                        "INSERT INTO ride_statuses (id, ride_id, status) VALUES (:id, :ride_id, :status)"
                    ),
                    {"id": str(ULID()), "ride_id": ride.id, "status": "CARRYING"},
                )
            case _:
                raise HTTPException(
                    status_code=status.HTTP_400_BAD_REQUEST, detail="invalid status"
                )
