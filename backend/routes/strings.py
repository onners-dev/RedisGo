from fastapi import APIRouter, HTTPException
from models import SetRequest
from redis_client import send_redis_command

router = APIRouter()

@router.post("/set")
def set_key(req: SetRequest):
    resp = send_redis_command(f"SET {req.key} {req.value}")
    if "+OK" in resp:
        return {"success": True}
    raise HTTPException(status_code=400, detail=resp)

@router.get("/get/{key}")
def get_key(key: str):
    resp = send_redis_command(f"GET {key}")
    if resp.startswith("$"):
        return {"value": resp.split('\r\n')[1]}
    return {"value": None}
