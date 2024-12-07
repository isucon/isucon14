import time
from http import HTTPStatus

from fastapi import APIRouter, Request, HTTPException
from sqlalchemy import text
from sse_starlette.sse import EventSourceResponse

from .app_handlers import (
    AppGetNotificationResponseChair,
    AppGetNotificationResponseData,
    Coordinate,
    calculate_discounted_fare,
    get_chair_stats,
    get_latest_ride_status,
)
from .middlewares import app_auth_middleware
from .models import Chair, Ride
from .sql import engine
from .utils import timestamp_millis

router = APIRouter()


@router.get("/api/app/notification")
def app_get_notification_sse(
    request: Request,
) -> EventSourceResponse:
    user = app_auth_middleware(request.cookies.get("app_session"))
    last_ride = None
    last_ride_status = None

    def event_stream():
        nonlocal request

        def f():
            nonlocal last_ride, last_ride_status, user
            with engine.begin() as conn:
                row = conn.execute(
                    text(
                        "SELECT * FROM rides WHERE user_id = :user_id ORDER BY created_at DESC LIMIT 1"
                    ),
                    {"user_id": user.id},
                ).fetchone()
                if row is None:
                    yield "data: {}"

                ride = Ride.model_validate(row)
                status = get_latest_ride_status(conn, ride.id)
                if (
                    (last_ride is not None)
                    and (ride.id == last_ride.id)
                    and (status == last_ride_status)
                ):
                    yield "data: {}"

                fare = calculate_discounted_fare(
                    conn,
                    user.id,
                    ride,
                    ride.pickup_latitude,
                    ride.pickup_longitude,
                    ride.destination_latitude,
                    ride.destination_longitude,
                )
                app_get_notification_response_chair = None

                if ride.chair_id:
                    row = conn.execute(
                        text("SELECT * FROM chairs WHERE id = :chair_id"),
                        {"chair_id": ride.chair_id},
                    ).fetchone()
                    if row is None:
                        raise HTTPException(
                            status_code=HTTPStatus.INTERNAL_SERVER_ERROR
                        )
                    chair = Chair.model_validate(row)

                    row = conn.execute(
                        text("SELECT * FROM chairs WHERE id = :chair_id"),
                        {"chair_id": ride.chair_id},
                    ).fetchone()
                    if row is None:
                        raise HTTPException(
                            status_code=HTTPStatus.INTERNAL_SERVER_ERROR
                        )
                    stats = get_chair_stats(conn, chair.id)
                    app_get_notification_response_chair = (
                        AppGetNotificationResponseChair(
                            id=chair.id,
                            name=chair.name,
                            model=chair.model,
                            stats=stats,
                        )
                    )
                data = AppGetNotificationResponseData(
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
                    chair=app_get_notification_response_chair,
                    created_at=timestamp_millis(ride.created_at),
                    updated_at=timestamp_millis(ride.updated_at),
                )
                yield str(data.model_dump_json(exclude_none=True))

        yield next(f())

        while True:
            # TODO: 同期的に止める方法を考える
            # if request.is_disconnected():
            #     break
            yield next(f())
            time.sleep(0.1)

    return EventSourceResponse(event_stream(), media_type="text/event-stream")
