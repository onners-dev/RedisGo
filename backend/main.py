import socket
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware

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

def send_redis_command(cmd: str) -> str:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect(("localhost", 6379))
        s.sendall((cmd + "\r\n").encode())
        data = s.recv(4096)
        return data.decode()

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
