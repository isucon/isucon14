import logging
import random
import string
import subprocess
from typing import Annotated

import sqlalchemy
from fastapi import FastAPI
from pydantic import BaseModel, StringConstraints
from sqlalchemy import text
from ulid import ULID

logging.basicConfig(level=logging.INFO)
logging.getLogger("sqlalchemy.engine").setLevel(logging.INFO)
app = FastAPI()
engine = sqlalchemy.create_engine("mysql+pymysql://isucon:isucon@localhost/isuride")
connection = engine.connect()


class PostInitializeRequest(BaseModel):
    payment_server: str


@app.post("/api/initialize")
def post_initialize() -> dict[str, str]:
    # TODO: fix path
    # TODO: golang output
    subprocess.run("../sql/init.sh", check=True)
    return {"language": "python"}


# owner_handlers.go
class PostOwnerRegisterRequest(BaseModel):
    name: Annotated[str, StringConstraints(min_length=1)]


class PostOwnerRegisterResponse(BaseModel):
    id: str


@app.post("/api/owner/register", status_code=201)
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


@app.post("/api/chair/register", status_code=201)
def chair_post_register():
    # TODO: implement
    # https://github.com/isucon/isucon14/blob/9571164b2b053f453dc0d24e0202d95c2fef253b/webapp/go/chair_handlers.go#L23

    chair_id = str(ULID())
    # TODO: should mimic secureRandomStr
    access_token = "".join(random.sample(string.ascii_letters + string.digits, 32))
    return {"access_token": access_token, "id": chair_id}


@app.post("/api/app/register", status_code=201)
def app_post_register():
    user_id = str(ULID())
    return {"id": user_id}


@app.get("/api/chair/notification", status_code=204)
def char_get_notification():
    # TODO: implement
    # https://github.com/isucon/isucon14/blob/9571164b2b053f453dc0d24e0202d95c2fef253b/webapp/go/chair_handlers.go#L141
    pass


@app.post("/api/app/payment-methods", status_code=204)
def app_post_payment_methods():
    # TODO: implement
    # https://github.com/isucon/isucon14/blob/9571164b2b053f453dc0d24e0202d95c2fef253b/webapp/go/app_handlers.go#L64
    pass


@app.post("/api/chair/activate", status_code=204)
def chair_post_activate():
    pass


@app.post("/api/chair/coordinate")
def chair_post_coordinate():
    return {"datetime": "2024-11-01T00:00:00Z"}  # RFC3339


@app.post("/api/app/requests", status_code=202)
def app_post_requests():
    request_id = str(ULID())
    return {"request_id": request_id}


@app.get("/api/app/notification", status_code=204)
def app_get_notification():
    pass


@app.get("/api/owner/chairs", status_code=200)
def owner_get_chairs():
    return {"chairs": []}


@app.get("/api/app/users", status_code=200)
def app_get_users():
    pass


@app.post("/api/app/users", status_code=200)
def app_post_users():
    pass
