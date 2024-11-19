"""
以下の移植
https://github.com/isucon/isucon14/blob/main/webapp/go/owner_handlers.go

TODO: このdocstringを消す
"""

import random
import string
from typing import Annotated

from fastapi import APIRouter
from pydantic import BaseModel, StringConstraints
from sqlalchemy import text
from ulid import ULID

from ..sql import engine

router = APIRouter(prefix="/api/owner")


class PostOwnerRegisterRequest(BaseModel):
    name: Annotated[str, StringConstraints(min_length=1)]


class PostOwnerRegisterResponse(BaseModel):
    id: str


@router.post("/register", status_code=201)
def owner_post_register(r: PostOwnerRegisterRequest) -> PostOwnerRegisterResponse:
    # TODO: implement
    # https://github.com/isucon/isucon14/blob/9571164b2b053f453dc0d24e0202d95c2fef253b/webapp/go/owner_handlers.go#L20

    owner_id = str(ULID())
    # TODO: should mimic secureRandomStr
    access_token = "".join(random.sample(string.ascii_letters + string.digits, 32))

    with engine.begin() as conn:
        conn.execute(
            text(
                "INSERT INTO owners (id, name, access_token) VALUES (:id, :name, :access_token)"
            ),
            {"id": owner_id, "name": r.name, "access_token": access_token},
        )
        return PostOwnerRegisterResponse(id=owner_id)


@router.get("/api/owner/chairs", status_code=200)
def owner_get_chairs():
    return {"chairs": []}
