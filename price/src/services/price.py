from grpc.aio import ServicerContext
from proto_gen.price_pb2_grpc import PriceServicer
from proto_gen.price_pb2 import Empty

class PriceService(PriceServicer):
    def __init__(self):
        pass

    async def HealthCheck(self, request: Empty, context: ServicerContext):
        return Empty()
