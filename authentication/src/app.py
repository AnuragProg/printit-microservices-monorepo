import asyncio
from src.client.postgres_client import PostgresClient
from src.client.redis_client import RedisClient

async def main():
    print('Application started')
    pg_client = PostgresClient()
    await pg_client.connect()
    redis_client = RedisClient()
    print('Application ended')

if __name__ == "__main__":
    asyncio.run(main())

