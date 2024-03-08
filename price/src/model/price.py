from client.postgres_client import PostgresClient


class PriceModel:

    def __init__(self, pg_client: PostgresClient):
        self.pg_client = pg_client

    async def create_table(self):
        await self.pg_client.execute(
            '''
            CREATE TABLE IF NOT EXISTS prices(
                _id UUID PRIMARY KEY,
                shop_id BYTEA NOT NULL,
                shopkeeper_id UUID NOT NULL,
                color TEXT CHECK(color in ('black-white', 'color')),
                page_size TEXT CHECK(page_size in ('A4')),
                single_sided_price DECIMAL NOT NULL CHECK(single_sided_price >= 0),
                double_sided_price DECIMAL NOT NULL CHECK(double_sided_price >= 0),
            );
            '''
        )
