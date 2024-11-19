"""
以下の移植
https://github.com/isucon/isucon14/blob/main/webapp/go/app_handlers.go

TODO: このdocstringを消す
"""

from fastapi import APIRouter
from ulid import ULID

router = APIRouter(prefix="/api/app")


@router.post("/register", status_code=201)
def app_post_register():
    user_id = str(ULID())
    return {"id": user_id}


@router.post("/payment-methods", status_code=204)
def app_post_payment_methods():
    # TODO: implement
    # https://github.com/isucon/isucon14/blob/9571164b2b053f453dc0d24e0202d95c2fef253b/webapp/go/app_handlers.go#L64
    pass


@router.post("/requests", status_code=202)
def app_post_requests():
    request_id = str(ULID())
    return {"request_id": request_id}


@router.get("/notification", status_code=204)
def app_get_notification():
    pass


@router.get("/users", status_code=200)
def app_get_users():
    pass


@router.post("/users", status_code=200)
def app_post_users():
    pass
