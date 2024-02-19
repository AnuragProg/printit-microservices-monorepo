import bcrypt


ROUNDS = 10


def hash(pw: str):
    salt = bcrypt.gensalt(rounds=ROUNDS)
    pw_b = bytes(pw, 'utf-8')
    hash = bcrypt.hashpw(pw_b, salt)
    return hash


def compare(pw: str, hash: str):
    pw_b, hash_b = bytes(pw, 'utf-8'), bytes(hash, 'utf-8')
    return bcrypt.checkpw(pw_b, hash_b)
