import uuid
import json
from enum import Enum
from src.client.postgres_client import PostgresClient
from src.client.redis_client import RedisClient
from src.util.secure_password_hasher import hash, compare


class UserType(Enum):
    CUSTOMER = 'customer'
    SHOPKEEPER = 'shopkeeper'

class User:
    def __init__(
        self,
        pg_client: PostgresClient,
        redis_client: RedisClient
    ):
        self.pg_client = pg_client
        self.redis_client = redis_client

    async def create_table(self):
        await self.pg_client.execute(
            '''
            CREATE TABLE IF NOT EXISTS users(
                _id UUID PRIMARY KEY,
                name TEXT,
                email TEXT UNIQUE,
                password TEXT,
                user_type TEXT CHECK(user_type IN ('customer', 'shopkeeper'))
            );
            '''
        )

    async def save(
        self,
        name: str,
        email: str,
        password: str,
        user_type: UserType
    ):
        '''
        Saves user data in Database
        NOTE: password is not hashed, and hence hashed password is to be provided
        REASON: commiting of user info is only done after verification, and hence verification process
        is to be done with hashed password
        '''
        _id = uuid.uuid4()
        await self.pg_client.execute(
            '''
            INSERT INTO users(
                _id,
                name,
                email,
                password,
                user_type
            )
            VALUES($1, $2, $3, $4, $5);
            '''
        , _id, name, email, password, user_type.name)


    async def cache(
        self,
        key: str,
        ttl: int,
        name: str,
        email: str,
        password: str,
        user_type: UserType
    ):
        '''
        Saves user data in Cache
        '''
        password = hash(password)
        val = json.dumps({
            name: name,
            email: email,
            password: password,
            user_type: user_type.name
        })
        await self.redis_client.set(key, ttl, val)







