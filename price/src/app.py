import os
import asyncio
import uvicorn
import grpc.aio
from client.postgres_client import PostgresClient
from fastapi import FastAPI

REST_HOST = os.getenv('REST_HOST', 'localhost')
REST_PORT = int(os.getenv('REST_PORT', '3000'))
GRPC_ADDR = os.getenv('GRPC_ADDR', '[::]:50051')


async def main():

    # setting up db
    pg_client = PostgresClient()
    await pg_client.connect()

    # rest app
    rest_app = FastAPI()

    # grpc app
    grpc_app = grpc.aio.server()
    grpc_app.add_insecure_port(GRPC_ADDR)
    #add_AuthenticationServicer_to_server(AuthenticationService(user_model), grpc_app)


    # setting up routers

    # include routers in main app

    # run the whole thing
    await grpc_app.start()
    print('GRPC server started...')
    print('REST server started...')
    config = uvicorn.Config(rest_app, host=REST_HOST, port=REST_PORT)
    server = uvicorn.Server(config)
    await server.serve()
    print('REST server stopped!')

    # gracefully close all services
    async def close_grpc():
        await grpc_app.stop(None)
        print('GRPC server stopped!')
    await asyncio.gather(
        close_grpc(),
        pg_client.close()
    )


if __name__ == "__main__":
    asyncio.run(main())

