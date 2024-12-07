import binascii
import os
from datetime import datetime, timedelta

from .models import Ride

INITIAL_FARE = 500
FARE_PER_DISTANCE = 100

EPOCH = datetime(1970, 1, 1)


def secure_random_str(b: int) -> str:
    random_bytes: bytes = os.urandom(b)
    return binascii.hexlify(random_bytes).decode("utf-8")


def timestamp_millis(dt: datetime) -> int:
    return (dt - EPOCH) // timedelta(milliseconds=1)


def datetime_fromtimestamp_millis(t: int) -> datetime:
    return EPOCH + timedelta(milliseconds=t)


def calculate_fare(
    pickup_latitude: int, pickup_longitude: int, dest_latitude: int, dest_longitude: int
) -> int:
    metered_fare = FARE_PER_DISTANCE * calculate_distance(
        pickup_latitude, pickup_longitude, dest_latitude, dest_longitude
    )
    return INITIAL_FARE + metered_fare


# マンハッタン距離を求める
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
