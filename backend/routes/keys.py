from fastapi import APIRouter
from redis_client import send_redis_command

router = APIRouter()

@router.get("/keys")
def get_keys():
    resp = send_redis_command("KEYS")
    keys = []
    if resp.startswith("*"):
        lines = resp.strip().split('\r\n')[1:]
        for i in range(1, len(lines), 2):
            keys.append(lines[i])
    return {"keys": keys}
