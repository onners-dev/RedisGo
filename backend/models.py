from pydantic import BaseModel

class SetRequest(BaseModel):
    key: str
    value: str

class CounterRequest(BaseModel):
    key: str
    action: str  # "incr" or "decr"

class CLIRequest(BaseModel):
    cmd: str

class HashSetRequest(BaseModel):
    key: str
    field: str
    value: str

class HashDelRequest(BaseModel):
    key: str
    field: str
