use crate::data::shop::Shop;
use mongodb::{Client, Collection, Database};
use std::sync::Mutex;


pub struct MongoManager{
    client: Client,
    printit_db: Database,
    pub shop_col: Collection<Shop>
}


impl MongoManager{
    pub fn new(client: Client) -> Self {
        let printit_db = client.database("prinit");
        let shop_col = printit_db.collection("shops");
        MongoManager{ client, printit_db, shop_col }
    }
}
