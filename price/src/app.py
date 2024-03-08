import os
import asyncio
import uvicorn
import grpc.aio
from fastapi import FastAPI
from model.price import PriceModel
from routes.price import PriceRouter
from services.price import PriceService
from client.postgres_client import PostgresClient
from proto_gen.auth_pb2_grpc import AuthenticationStub
from proto_gen.price_pb2_grpc import add_PriceServicer_to_server


REST_HOST = os.getenv('REST_HOST', 'localhost')
REST_PORT = os.getenv('REST_PORT', '3003')
GRPC_ADDR = os.getenv('GRPC_ADDR', '[::]:50054')
AUTH_GRPC_URI = os.getenv('AUTH_GRPC_URI', 'localhost:50051')


async def main():

    # connecting to the database
    pg_client = PostgresClient()
    await pg_client.connect()

    # connecting to grpc services
    auth_grpc_channel = grpc.aio.insecure_channel(AUTH_GRPC_URI)
    auth_grpc_client = AuthenticationStub(auth_grpc_channel)

    # setting up data models
    price_model = PriceModel(pg_client=pg_client)
    await price_model.create_table()

    # rest app
    rest_app = FastAPI()

    # grpc app
    grpc_app = grpc.aio.server()
    grpc_app.add_insecure_port(GRPC_ADDR)
    add_PriceServicer_to_server(PriceService(), grpc_app)

    # setting up routers
    price_router = PriceRouter(auth_grpc_client=auth_grpc_client)
    price_router.setup_routes()

    # include routers in main app
    rest_app.include_router(price_router.router, prefix='/price')

    # run grpc + rest
    await grpc_app.start()
    print('GRPC server started...')

    print('REST server started...')
    config = uvicorn.Config(rest_app, host=REST_HOST, port=int(REST_PORT))
    server = uvicorn.Server(config)
    await server.serve()
    print('REST server stopped!')

    # gracefully close (rest + grpc)
    async def close_grpc():
        await grpc_app.stop(None)
        print('GRPC server stopped!')
    await asyncio.gather(
        close_grpc(),
        pg_client.close()
    )

if __name__ == "__main__":
    asyncio.run(main())

