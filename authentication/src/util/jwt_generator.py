import jwt

SECRET = "secret"
ALGO = "HS256"

def decode_jwt(token: str):
    return jwt.decode(token, SECRET, algorithms=[ALGO])


def generate_jwt(**kwargs):
    return jwt.encode(kwargs, SECRET, algorithm=ALGO)

