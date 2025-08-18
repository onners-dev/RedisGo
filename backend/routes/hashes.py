from fastapi import APIRouter, HTTPException
from models import HashSetRequest, HashDelRequest
from redis_client import send_redis_command

router = APIRouter()

@router.post("/hash/set")
def set_hash_field(req: HashSetRequest):
    resp = send_redis_command(f"HSET {req.key} {req.field} {req.value}")
    if resp.startswith(":"):
        return {"created": bool(int(resp[1:].strip()))}
    raise HTTPException(status_code=400, detail=resp)

@router.get("/hash/get")
def get_hash_field(key: str, field: str):
    resp = send_redis_command(f"HGET {key} {field}")
    if resp.startswith("$"):
        lines = resp.strip().split('\r\n')
        if len(lines) >= 2:
            return {"value": lines[1]}
        else:
            return {"value": None}
    return {"value": None}

@router.get("/hash/all")
def get_hash_all(key: str):
    resp = send_redis_command(f"HGETALL {key}")
    if not resp.startswith("*"):
        return {"fields": {}}
    lines = resp.strip().split('\r\n')[1:]
    out = {}
    for i in range(0, len(lines), 2):
        if i + 1 < len(lines):
            out[lines[i]] = lines[i + 1]
    return {"fields": out}

@router.post("/hash/del")
def del_hash_field(req: HashDelRequest):
    resp = send_redis_command(f"HDEL {req.key} {req.field}")
    if resp.startswith(":"):
        return {"deleted": bool(int(resp[1:].strip()))}
    raise HTTPException(status_code=400, detail=resp)
