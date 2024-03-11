import os
import grpc.aio
from proto_gen.auth_pb2_grpc import AuthenticationStub


AUTH_GRPC_URI = os.getenv('AUTH_GRPC_URI', 'localhost:50051')

__auth_grpc_client = None
__auth_grpc_channel = None

async def get_auth_grpc_client():
    global __auth_grpc_client
    global __auth_grpc_channel

    if __auth_grpc_client is None:
        __auth_grpc_channel = grpc.aio.insecure_channel(AUTH_GRPC_URI)
        __auth_grpc_client = AuthenticationStub(__auth_grpc_channel)

    return __auth_grpc_client


async def close_auth_grpc_client():
    if __auth_grpc_channel is not None:
        await __auth_grpc_channel.close()
