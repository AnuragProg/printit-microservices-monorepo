from typing import Union
import asyncpg

class PostgresClient:
    def __init__(self):
        self.pool: Union[asyncpg.Pool, None] = None

    async def connect(self):
        self.pool = await asyncpg.create_pool(
            user="root",
            password="root",
            database="users",
            host="localhost",
            max_size=100
        )
        print('Connected to Postgres')

    async def fetch(self, query, *args):
        if not self.pool:
            raise Exception('pool not initialized')

        async with self.pool.acquire() as conn:
            if len(args)==0:
                return await conn.fetch(query)
            return await conn.fetch(query, *args)


    async def execute(self, query, *args):
        if not self.pool:
            raise Exception('pool not initialized')
        async with self.pool.acquire() as conn:
            if len(args)==0:
                return await conn.execute(query)
            return await conn.execute(query, *args)
