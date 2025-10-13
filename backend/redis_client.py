import socket
import os

REDIS_HOST = os.environ.get("REDIS_HOST", "redisgo")
REDIS_PORT = int(os.environ.get("REDIS_PORT", 6379))

def send_redis_command(cmd: str) -> str:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((REDIS_HOST, REDIS_PORT))
        s.sendall((cmd + "\r\n").encode())
        data = s.recv(4096)
        return data.decode()
