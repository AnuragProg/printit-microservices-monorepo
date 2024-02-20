import uuid
import json
from enum import Enum
from src.errors.not_found import NotFound
from src.client.postgres_client import PostgresClient
from src.client.redis_client import RedisClient
from src.util.secure_password_hasher import hash, compare


class User:
    def __init__(self, _id: str, name: str, email: str, user_type: str):
        self._id = _id
        self.name = name
        self.email = email
        self.user_type = user_type

class UserType(Enum):
    CUSTOMER = 'customer'
    SHOPKEEPER = 'shopkeeper'

    @staticmethod
    def get(user_type: str):
        types = {
            UserType.CUSTOMER.value: UserType.CUSTOMER,
            UserType.SHOPKEEPER.value: UserType.SHOPKEEPER
        }
        return types[user_type]

class Purpose(Enum):
    FORGOT_PASSWORD = 'forgot-password'
    SIGNUP = 'signup'

class UserModel:
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

    async def get_user_from_db(
        self,
        email: str,
        password: str
    ):
        try:
            password = hash(password)
            res = await self.pg_client.fetch(
                '''
                SELECT * FROM users WHERE email=$1 AND password=$2 LIMIT 1;
                '''
            , email, password)
            row = res[0]
            return User(row['_id'], row['name'], row['email'], row['user_type'])
        except Exception:
            raise NotFound()

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
        user_type: UserType,
        purpose: Purpose
    ):
        '''
        Saves user data in Cache
        '''
        password = hash(password)
        val = json.dumps({
            'name': name,
            'email': email,
            'password': password,
            'user_type': user_type.value,
            'purpose': purpose.value
        })
        await self.redis_client.set(key, ttl, val)


    async def get_otp_info(self, key: str):
        val = await self.redis_client.get(key)
        if val is None:
            raise NotFound()
        val_str = json.loads(val)
        return val_str

    async def del_otp_info(self, key: str):
        await self.redis_client.delete(key)




