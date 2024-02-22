import jwt
import datetime

SECRET = "secret"
ALGO = "HS256"
TOKEN_EXPIRY = 30 # days

def decode_jwt(token: str):
    return jwt.decode(token, SECRET, algorithms=[ALGO])


def generate_jwt(**kwargs):
    payload = {
        **kwargs,
        'exp': datetime.datetime.utcnow() + datetime.timedelta(days=TOKEN_EXPIRY)
    }
    return jwt.encode(payload, SECRET, algorithm=ALGO)

