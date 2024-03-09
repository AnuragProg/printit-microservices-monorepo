tonic::include_proto!("shop_grpc");
use shop_server::Shop;
use tonic::{Request, Response, Status};
use crate::data::shop::Shop as ShopDoc;
use crate::data::location::Location as LocationDoc;
use crate::client::mongo_client::MongoManager;
use mongodb::bson::{oid::ObjectId, doc};
use std::sync::Arc;
use std::convert::From;



impl From<LocationDoc> for Location{
    fn from(location: LocationDoc) -> Self {
        Location {
            lng: location.coordinates[0],
            lat: location.coordinates[1]
        }
    }
}

impl From<ShopDoc> for ShopInfo{
    fn from(val: ShopDoc) -> Self {
        ShopInfo{
            id: val._id.to_hex(),
            user_id: val.user_id,
            name: val.name,
            contact: val.contact,
            email: val.email,
            location: Some(Location::from(val.location)),
            created_at: val.created_at.try_to_rfc3339_string().unwrap(),
            updated_at: val.updated_at.try_to_rfc3339_string().unwrap(),
        }
    }
}


#[derive(Debug)]
pub struct ShopService{
    mongo_manager: Arc<MongoManager>
}

impl ShopService{
    pub fn new(mongo_manager: Arc<MongoManager>) -> Self {
        ShopService{ mongo_manager }
    }
}


#[tonic::async_trait]
impl Shop for ShopService{
    async fn health_check(&self, _req: Request<Empty>) ->  Result<Response<Empty>, Status>{
        Ok(Response::new(Empty{}))
    }
    async fn verify_shop(&self, req: Request<ShopId>) -> Result<Response<ShopInfo>, Status>{
        let shop_id = ObjectId::parse_str(req.into_inner().id);
        if shop_id.is_err() {
            return Err(Status::invalid_argument("invalid shop id"));
        }
        let shop_result = self.mongo_manager.shop_col.find_one(doc!{"_id": shop_id.unwrap()}, None).await;
        if let Err(err) = shop_result {
            return Err(Status::internal(err.to_string()));
        }
        let shop = shop_result.unwrap();
        if shop.is_none() {
            return Err(Status::not_found("shop not found"));
        }
        Ok(Response::new(shop.unwrap().into()))
    }
}
