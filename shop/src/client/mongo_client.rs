use crate::data::shop::Shop;
use mongodb::{Client, Collection, Database, IndexModel, bson::doc};


pub struct MongoManager{
    client: Client,
    printit_db: Database,
    pub shop_col: Collection<Shop>
}

impl MongoManager{
    pub async fn new() -> Self {
        // creating client
        let client = Client::with_uri_str("mongodb://localhost:27017/?maxPoolSize=100").await.unwrap();

        // db
        let printit_db = client.database("printit");

        // setting shop collection
        let shop_col = printit_db.collection("shops");
        let shop_geo_index = IndexModel::builder()
            .keys(doc!{"location": "2dsphere"})
            .build();
        shop_col.create_index(shop_geo_index, None).await.unwrap();

        MongoManager{ client, printit_db, shop_col }
    }
}
