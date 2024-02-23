import grpc
from proto_gen.auth_pb2_grpc import AuthenticationServicer
from proto_gen.auth_pb2 import Token, User, UserType as GrpcUserType
from grpc.aio import ServicerContext
from util.jwt_generator import decode_jwt
from model.user import UserModel, UserType as LocalUserType


class AuthenticationService(AuthenticationServicer):
    def __init__(self, user_model: UserModel):
        self.user_model = user_model

    async def VerifyToken(self, request: Token, context: ServicerContext):
        try:
            id = decode_jwt(token=request.token)['_id']
            user = await self.user_model.get_user_from_db_by_id(id)
            user_type = GrpcUserType.UNDEFINED
            if user.user_type == LocalUserType.CUSTOMER.value:
                user_type = GrpcUserType.CUSTOMER
            elif user.user_type == LocalUserType.SHOPKEEPER.value:
                user_type = GrpcUserType.SHOPKEEPER
            return User(
                _id=id,
                name=user.name,
                email=user.email,
                user_type=user_type
            )
        except Exception:
            context.set_code(grpc.StatusCode.UNAUTHENTICATED)
            context.set_details('Invalid token')
            return None
