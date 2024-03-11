import os
import asyncio
import uvicorn
import grpc.aio
from fastapi import FastAPI
from client.auth_grpc_client import close_auth_grpc_client
from routes.price import router as price_router
from services.price import PriceService
from client.postgres_client import close_pg_client
from proto_gen.price_pb2_grpc import add_PriceServicer_to_server


REST_HOST = os.getenv('REST_HOST', 'localhost')
REST_PORT = os.getenv('REST_PORT', '3003')
GRPC_ADDR = os.getenv('GRPC_ADDR', '[::]:50054')


async def main():

    # rest app
    rest_app = FastAPI()

    # include routers in main app
    rest_app.include_router(price_router, prefix='/price')

    # grpc app
    grpc_app = grpc.aio.server()
    grpc_app.add_insecure_port(GRPC_ADDR)
    add_PriceServicer_to_server(PriceService(), grpc_app)

    # run grpc + rest
    await grpc_app.start()
    print('GRPC server started...')

    print('REST server started...')
    config = uvicorn.Config(
        rest_app,
        host=REST_HOST,
        port=int(REST_PORT),
    )
    server = uvicorn.Server(config)
    await server.serve()
    print('REST server stopped!')

    # gracefully close (rest + grpc)
    async def close_grpc():
        await grpc_app.stop(None)
        print('GRPC server stopped!')

    await asyncio.gather(
        close_grpc(),               # closing grpc server
        close_pg_client(),          # closing pg client
        close_auth_grpc_client(),   # closing auth grpc client
    )

if __name__ == "__main__":
    asyncio.run(main())

