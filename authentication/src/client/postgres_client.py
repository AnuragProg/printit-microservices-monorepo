import asyncpg


class PostgresClient:
    def __init__(self):
        self.pool = None

    async def connect(self):
        self.pool = await asyncpg.create_pool(
            user="root",
            password="root",
            database="users",
            host="localhost"
        )

    async def execute(self, query, *args):
        if not self.pool:
            raise Exception('pool not initialized')
        async with self.pool.acquire() as conn:
            if(len(args)==0):
                return await conn.execute(query)
            return await conn.execute(query, *args)
