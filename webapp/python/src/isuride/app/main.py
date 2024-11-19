import logging
import subprocess

from fastapi import FastAPI
from pydantic import BaseModel

from .routers import apps, chairs, owners

logging.basicConfig(level=logging.INFO)
logging.getLogger("sqlalchemy.engine").setLevel(logging.INFO)
app = FastAPI()
app.include_router(apps.router)
app.include_router(chairs.router)
app.include_router(owners.router)


class PostInitializeRequest(BaseModel):
    payment_server: str


@app.post("/api/initialize")
def post_initialize() -> dict[str, str]:
    # TODO: fix path
    # TODO: golang output
    subprocess.run("../sql/init.sh", check=True)
    return {"language": "python"}
