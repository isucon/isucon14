import binascii
import os
from datetime import datetime, timedelta
from zoneinfo import ZoneInfo

# if typing.TYPE_CHECKINGの中に入れた方がいいかも
from .models import Ride

INITIAL_FARE: int = 500
FARE_PER_DISTANCE: int = 100
UTC: ZoneInfo = ZoneInfo("UTC")

EPOCH = datetime(1970, 1, 1, tzinfo=UTC)


def secure_random_str(b: int) -> str:
    random_bytes: bytes = os.urandom(b)
    return binascii.hexlify(random_bytes).decode("utf-8")


def timestamp_millis(dt: datetime) -> int:
    dt = dt.astimezone(UTC)
    return (dt - EPOCH) // timedelta(milliseconds=1)


def datetime_fromtimestamp_millis(t: int) -> datetime:
    return EPOCH + timedelta(milliseconds=t)


def calculate_fare(pickup_latitude, pickup_longitude, dest_latitude, dest_longitude):
    metered_fare = FARE_PER_DISTANCE * calculate_distance(
        pickup_latitude, pickup_longitude, dest_latitude, dest_longitude
    )
    return INITIAL_FARE + metered_fare


def calculate_distance(
    a_latitude: int, a_longitude: int, b_latitude: int, b_longitude: int
) -> int:
    return abs(a_latitude - b_latitude) + abs(a_longitude - b_longitude)


def calculate_sale(ride: Ride) -> int:
    return calculate_fare(
        ride.pickup_latitude,
        ride.pickup_longitude,
        ride.destination_latitude,
        ride.destination_longitude,
    )


def sum_sales(rides: list[Ride]) -> int:
    sale = 0
    for ride in rides:
        sale += calculate_sale(ride)
    return sale