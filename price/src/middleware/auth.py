from typing import Annotated

from fastapi import HTTPException, Header
from proto_gen.auth_pb2 import Token, UserType, User
from client.auth_grpc_client import get_auth_grpc_client


async def auth_user(authorization: Annotated[str|None, Header()]):
    if authorization is None:
        raise HTTPException(401, detail='missing auth token')
    split = authorization.split(' ')
    if len(split) < 2:
        raise HTTPException(401, detail='missing auth token')
    token = split[1]
    try:
        auth_grpc_client = await get_auth_grpc_client()
        user = await auth_grpc_client.VerifyToken(Token(token=token))
        return user
    except Exception as e:
        raise HTTPException(401, detail=str(e))


async def auth_shopkeeper(authorization: Annotated[str|None, Header()]):
    if authorization is None:
        raise HTTPException(401, detail='missing auth token')
    split = authorization.split(' ')
    if len(split) < 2:
        raise HTTPException(401, detail='missing auth token')
    token = split[1]
    try:
        auth_grpc_client = await get_auth_grpc_client()
        user: User = await auth_grpc_client.VerifyToken(Token(token=token))
        if user.user_type != UserType.SHOPKEEPER:
            raise Exception("user not a shopkeeper")
        return user
    except Exception as e:
        raise HTTPException(401, detail=str(e))

