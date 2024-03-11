import uuid
from enum import Enum
from typing import List, Optional
from asyncpg import Record
from bson import ObjectId as oid
from client.postgres_client import PostgresClient, get_pg_client
from errors.not_found import NotFound


# price and product are being used interchangeably

class PricePageSize(Enum):
    A4 = 'A4'

    @staticmethod
    def get(elem: str):
        for page_size in PricePageSize:
            if page_size.value == elem:
                return page_size
        raise ValueError()


class PriceColor(Enum):
    BLACK_WHITE = 'black-white'
    COLOR       = 'color'

    @staticmethod
    def get(elem: str):
        for color in PriceColor:
            if color.value == elem:
                return color
        raise ValueError()


class Price:
    def __init__(
        self,
        shop_id: oid,
        shopkeeper_id: uuid.UUID,
        color: PriceColor,
        single_sided_price: float,
        double_sided_price: float,
        _id: uuid.UUID | None = None,
        page_size: PricePageSize | None = None
    ):
        if _id is None:
            self._id = uuid.uuid4()
        else:
            self._id = _id
        self.shop_id = shop_id
        self.shopkeeper_id = shopkeeper_id
        self.color = color
        if page_size == None:
            self.page_size = PricePageSize.A4
        else:
            self.page_size = page_size
        self.single_sided_price = single_sided_price
        self.double_sided_price = double_sided_price

    def to_dict(self):
        return {
            "_id": self._id.hex,
            "shop_id": str(self.shop_id),
            "shopkeeper_id": self.shopkeeper_id.hex,
            "color": self.color.value,
            "single_sided_price": self.single_sided_price,
            "double_sided_price": self.double_sided_price,
            "page_size": self.page_size.value
        }
    @staticmethod
    def pg_record_to_price(record):
        return Price(
            shop_id = oid(record['shop_id']),
            shopkeeper_id = record['shopkeeper_id'],
            color = PriceColor.get(record['color']),
            single_sided_price = record['single_sided_price'],
            double_sided_price = record['double_sided_price'],
            _id = record['_id'],
            page_size=PricePageSize.get(record['page_size'])
        )

class PriceModel:

    def __init__(self, pg_client: PostgresClient):
        self.pg_client = pg_client

    async def _create_table(self):
        await self.pg_client.execute(
            '''
            CREATE TABLE IF NOT EXISTS prices(
                _id UUID PRIMARY KEY,
                shop_id BYTEA NOT NULL,
                shopkeeper_id UUID NOT NULL,

                color TEXT CHECK(color in ('black-white', 'color')),
                page_size TEXT CHECK(page_size in ('A4')),

                single_sided_price FLOAT NOT NULL CHECK(single_sided_price >= 0),
                double_sided_price FLOAT NOT NULL CHECK(double_sided_price >= 0)
            );
            '''
        )

    async def insert_price(self, price: Price):
        await self.pg_client.execute(
            '''
            INSERT INTO prices(_id, shop_id, shopkeeper_id, color, page_size, single_sided_price, double_sided_price)
            VALUES($1, $2, $3, $4, $5, $6, $7);
            ''',
            price._id, price.shop_id.binary, price.shopkeeper_id, price.color.value, price.page_size.value, price.single_sided_price, price.double_sided_price
        )

    async def update_price(
        self,
        price_id: uuid.UUID,
        shop_id: oid,
        single_sided_price: Optional[float],
        double_sided_price: Optional[float],
        color: Optional[PriceColor],
        page_size: Optional[PricePageSize]
    ):
        query = 'UPDATE prices SET '
        param_id= 1
        clauses = []
        params  = []

        if single_sided_price is not None:
            clauses.append(f'single_sided_price = ${param_id}')
            params.append(single_sided_price)
            param_id += 1

        if double_sided_price is not None:
            clauses.append(f'double_sided_price = ${param_id}')
            params.append(double_sided_price)
            param_id += 1

        if color is not None:
            clauses.append(f'color = ${param_id}')
            params.append(color.value)
            param_id += 1

        if page_size is not None:
            clauses.append(f'page_size = ${param_id}')
            params.append(page_size.value)
            param_id += 1

        if param_id == 1:
            # no fields are given for updation
            return

        query += ', '.join(clauses)
        query += f' WHERE _id = ${param_id} AND shop_id = ${param_id + 1};'
        params.append(price_id)
        params.append(shop_id.binary)

        print(f'query is = {query}')
        print(f'params are = {params}')

        await self.pg_client.execute(query, *params)

    async def get_price_by_id(self, shop_id: oid, price_id: uuid.UUID):
        price_records = await self.pg_client.fetch(
            '''
            SELECT * FROM prices WHERE _id = $1 AND shop_id = $2;
            ''',
            price_id, shop_id.binary
        )
        if len(price_records) < 1:
            raise NotFound(msg='price not found', resource='price')
        return Price.pg_record_to_price(price_records[0])

    async def get_prices_by_shop_id(self, shop_id: oid):
        price_records = await self.pg_client.fetch(
            '''
            SELECT * FROM prices WHERE shop_id = $1;
            ''',
            shop_id.binary
        )
        return list(map(Price.pg_record_to_price, price_records))

    async def delete_price(self, price_id: uuid.UUID, shop_id: oid):
        await self.pg_client.execute(
            '''
            DELETE FROM prices WHERE _id = $1 AND shop_id = $2;
            ''',
            price_id, shop_id.binary
        )


_price_model = None

async def get_price_model():
    global _price_model
    if _price_model is None:
        pg_client = await get_pg_client()
        _price_model = PriceModel(pg_client)
        await _price_model._create_table()
    return _price_model
