from fastapi import FastAPI, HTTPException, Request
from pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware
import socket

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

class SetRequest(BaseModel):
    key: str
    value: str

class CounterRequest(BaseModel):
    key: str
    action: str  # "incr" or "decr"


class CLIRequest(BaseModel):
    cmd: str

def send_redis_command(cmd: str) -> str:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect(("localhost", 6379))
        s.sendall((cmd + "\r\n").encode())
        data = s.recv(4096)
        return data.decode()

@app.post("/cli")
def cli_command(req: CLIRequest):
    resp = send_redis_command(req.cmd)
    return {"resp": resp.strip()}

@app.post("/set")
def set_key(req: SetRequest):
    resp = send_redis_command(f"SET {req.key} {req.value}")
    if "+OK" in resp:
        return {"success": True}
    raise HTTPException(status_code=400, detail=resp)

@app.get("/get/{key}")
def get_key(key: str):
    resp = send_redis_command(f"GET {key}")
    if resp.startswith("$"):
        return {"value": resp.split('\r\n')[1]}
    return {"value": None}


@app.post("/counter")
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


@app.get("/keys")
def get_keys():
    resp = send_redis_command("KEYS")
    keys = []
    if resp.startswith("*"):
        lines = resp.strip().split('\r\n')[1:]
        for i in range(1, len(lines), 2):
            keys.append(lines[i])
    return {"keys": keys}
