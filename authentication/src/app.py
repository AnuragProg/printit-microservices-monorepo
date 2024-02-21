import asyncio
import uvicorn
from src.model.user import UserModel
from src.client.postgres_client import PostgresClient
from src.client.redis_client import RedisClient
from fastapi import FastAPI
from src.routes.user import UserRouter


async def main():
    # rest app
    rest = FastAPI()

    # for running async operations
    #event_loop = asyncio.get_event_loop()

    # setting up db and cache clients
    pg_client = PostgresClient()
    await pg_client.connect()
    redis_client = RedisClient()
    user_model = UserModel(pg_client, redis_client)
    await user_model.create_table()

    # setting up routers
    user_router = UserRouter(user_model)
    user_router.setup_routes()

    # include routers in main app
    rest.include_router(router=user_router.router, prefix='/user')

    config = uvicorn.Config(rest, host='localhost', port=3000)
    server = uvicorn.Server(config)

    # run the whole thing
    await server.serve()

if __name__ == "__main__":
    asyncio.run(main())

