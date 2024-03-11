import os
import grpc.aio
from proto_gen.shop_pb2_grpc import ShopStub


SHOP_GRPC_URI = os.getenv('SHOP_GRPC_URI', 'localhost:50053')

__shop_grpc_client = None
__shop_grpc_channel = None

async def get_shop_grpc_client():
    global __shop_grpc_client
    global __shop_grpc_channel

    if __shop_grpc_client is None:
        __shop_grpc_channel = grpc.aio.insecure_channel(SHOP_GRPC_URI)
        __shop_grpc_client = ShopStub(__shop_grpc_channel)

    return __shop_grpc_client


async def close_shop_grpc_client():
    if __shop_grpc_channel is not None:
        await __shop_grpc_channel.close()
