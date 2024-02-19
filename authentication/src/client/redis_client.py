import redis


CONNECTION_TIMEOUT = 60 # secs

class RedisClient:
    def __init__(self):
        self.pool = redis.BlockingConnectionPool(
            max_connections=100,
            timeout=CONNECTION_TIMEOUT,
            host='localhost',
            port=6379,
            db=0
        )
        self.client = redis.Redis(connection_pool=self.pool)

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
        with self.client.pipeline() as pipe:
            if ttl==0:
                # setting cache forever
                pipe.set(name=key, value=value)
            else:
                # setting cache for ttl
                pipe.setex(name=key, time=ttl, value=value)

    async def get(
        self,
        key: str
    ):
        '''
        Gets value of provided key
        '''
        with self.client.pipeline() as pipe:
            return pipe.get(key)




