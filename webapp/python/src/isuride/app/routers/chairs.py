"""
以下の移植
https://github.com/isucon/isucon14/blob/main/webapp/go/chair_handlers.go

TODO: このdocstringを消す
"""

import random
import string

from fastapi import APIRouter
from ulid import ULID

router = APIRouter(prefix="/api/chair")


@router.post("/register", status_code=201)
def chair_post_register():
    # TODO: implement
    # https://github.com/isucon/isucon14/blob/9571164b2b053f453dc0d24e0202d95c2fef253b/webapp/go/chair_handlers.go#L23

    chair_id = str(ULID())
    # TODO: should mimic secureRandomStr
    access_token = "".join(random.sample(string.ascii_letters + string.digits, 32))
    return {"access_token": access_token, "id": chair_id}


@router.get("/notification", status_code=204)
def char_get_notification():
    # TODO: implement
    # https://github.com/isucon/isucon14/blob/9571164b2b053f453dc0d24e0202d95c2fef253b/webapp/go/chair_handlers.go#L141
    pass


@router.post("/activate", status_code=204)
def chair_post_activate():
    pass


@router.post("/coordinate")
def chair_post_coordinate():
    return {"datetime": "2024-11-01T00:00:00Z"}  # RFC3339
