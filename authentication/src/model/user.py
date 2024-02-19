import uuid
from enum import Enum
from src.client.postgres_client import PostgresClient
from src.util.secure_password_hasher import hash, compare


class UserType(Enum):
    CUSTOMER = 'customer'
    SHOPKEEPER = 'shopkeeper'


class User:
    def __init__(self, postgres_client: PostgresClient):
        self.postgres_client = postgres_client

    async def create_table(self):
        await self.postgres_client.execute(
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
        '''
        _id = uuid.uuid4()
        await self.postgres_client.execute(
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
        pass

