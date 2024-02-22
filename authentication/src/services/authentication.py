import grpc
from src.proto_gen.auth_pb2_grpc import AuthenticationServicer
from src.proto_gen.auth_pb2 import Token, User
from grpc.aio import ServicerContext
from src.util.jwt_generator import decode_jwt
from src.model.user import UserModel


class AuthenticationService(AuthenticationServicer):
    def __init__(self, user_model: UserModel):
        self.user_model = user_model

    async def VerifyToken(self, request: Token, context: ServicerContext):
        try:
            id = decode_jwt(token=request.token)['_id']
            user = await self.user_model.get_user_from_db_by_id(id)
            return User(
                _id=id,
                name=user.name,
                email=user.email,
                user_type=user.user_type
            )
        except Exception:
            context.set_code(grpc.StatusCode.UNAUTHENTICATED)
            context.set_details('Invalid token')
            return None