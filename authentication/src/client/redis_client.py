import os
import redis.asyncio as redis

CONNECTION_TIMEOUT = 60 # secs

REDIS_URL = os.getenv('REDIS_URL', 'redis://localhost/0')

class RedisClient:
    def __init__(self):
        self.pool = redis.BlockingConnectionPool.from_url(
            REDIS_URL,
            max_connections=110,
            timeout=30
        )
        print('Connected to Redis')

    async def close(self):
        print('Closing Redis pool...')
        await self.pool.disconnect()
        print('Closed Redis pool...')

    async def set(
        self,
        key: str,
        ttl: int,
        value: str
    ):
        '''
        Sets value for key
        if ttl == 0: then cache is set forever
        else the cache is set with given ttl
        '''
        conn = redis.Redis(connection_pool=self.pool)
        if ttl==0:
            # setting cache forever
            await conn.set(name=key, value=value)
        else:
            # setting cache for ttl
            await conn.setex(name=key, time=ttl, value=value)
        await conn.aclose()

    async def get(self, key: str):
        '''
        Gets value of provided key
        '''
        conn = redis.Redis(connection_pool=self.pool)
        val = await conn.get(key)
        await conn.aclose()
        return val

    async def delete(self, key: str):
        conn = redis.Redis(connection_pool=self.pool)
        await conn.delete(key)
        await conn.aclose()



