from fastapi import APIRouter, HTTPException
from models import CounterRequest
from redis_client import send_redis_command

router = APIRouter()

@router.post("/counter")
def update_counter(req: CounterRequest):
    if req.action == "incr":
        resp = send_redis_command(f"INCR {req.key}")
    elif req.action == "decr":
        resp = send_redis_command(f"DECR {req.key}")
    else:
        raise HTTPException(status_code=400, detail="Invalid action for counter.")

    if resp.startswith(":"):
        return {"value": int(resp[1:].strip())}
    raise HTTPException(status_code=400, detail="Redis error: " + resp)
