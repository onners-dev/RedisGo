from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from routes.strings import router as strings_router
from routes.counters import router as counters_router
from routes.keys import router as keys_router
from routes.hashes import router as hashes_router
from models import CLIRequest
from redis_client import send_redis_command
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.post("/cli")
def cli_command(req: CLIRequest):
    resp = send_redis_command(req.cmd)
    return {"resp": resp.strip()}

# Mount routers with appropriate prefixes if you want, or just mount as-is
app.include_router(strings_router)
app.include_router(counters_router)
app.include_router(keys_router)
app.include_router(hashes_router)
