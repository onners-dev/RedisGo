import socket

def send_redis_command(cmd: str) -> str:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect(("localhost", 6379))
        s.sendall((cmd + "\r\n").encode())
        data = s.recv(4096)
        return data.decode()
