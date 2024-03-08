
from fastapi import APIRouter
from proto_gen.auth_pb2_grpc import AuthenticationStub




class PriceRouter:
    def __init__(self, auth_grpc_client: AuthenticationStub):
        self.router = APIRouter()
        self.auth_grpc_client = auth_grpc_client


    def setup_routes(self):
        pass
