import os
import asyncio
import uvicorn
import grpc.aio
from model.user import UserModel
from client.postgres_client import PostgresClient
from client.redis_client import RedisClient
from fastapi import FastAPI
from routes.user import UserRouter
from proto_gen.auth_pb2_grpc import add_AuthenticationServicer_to_server
from services.authentication import AuthenticationService

REST_HOST = os.getenv('REST_HOST', 'localhost')
REST_PORT = int(os.getenv('REST_PORT', '3000'))
GRPC_ADDR = os.getenv('GRPC_ADDR', '[::]:50051')


async def main():

    # setting up db
    pg_client = PostgresClient()
    await pg_client.connect()

    # setting up cache
    redis_client = RedisClient()
    user_model = UserModel(pg_client, redis_client)
    await user_model.create_table()


    # rest app
    rest_app = FastAPI()

    # grpc app
    grpc_app = grpc.aio.server()
    grpc_app.add_insecure_port(GRPC_ADDR)
    add_AuthenticationServicer_to_server(AuthenticationService(user_model), grpc_app)


    # setting up routers
    user_router = UserRouter(user_model)
    user_router.setup_routes()

    # include routers in main app
    rest_app.include_router(router=user_router.router, prefix='/user')

    config = uvicorn.Config(rest_app, host=REST_HOST, port=REST_PORT)
    server = uvicorn.Server(config)

    # run the whole thing
    await grpc_app.start()
    print('GRPC server started...')
    print('REST server started...')
    await server.serve()
    print('REST server stopped!')

    await grpc_app.stop(None)
    print('GRPC server stopped!')


if __name__ == "__main__":
    asyncio.run(main())

