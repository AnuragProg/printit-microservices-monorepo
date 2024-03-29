tonic::include_proto!("shop_grpc");
use shop_server::Shop;
use tonic::{Request, Response, Status};
use crate::client::mongo_client::MongoManager;
use mongodb::bson::{doc, oid::ObjectId};
use std::sync::Arc;



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
    async fn health_check(&self, _: Request<Empty>) ->  Result<Response<Empty>, Status>{
        Ok(Response::new(Empty{}))
    }
    async fn get_shop_by_id(&self, shop_id: Request<ShopId>) ->  Result<Response<ShopInfo>, Status>{
        let id = ObjectId::parse_str(shop_id.into_inner().id);
        if id.is_err() {
            return Err(Status::invalid_argument("invalid shop id"));
        }
        let shop_res = self.mongo_manager.shop_col.find_one(doc!{"_id": id.unwrap()}, None).await;
        if shop_res.is_err() {
            return Err(Status::internal("unable to access shop db"));
        }
        let shop_opt = shop_res.unwrap();
        if shop_opt.is_none() {
            return Err(Status::not_found("shop not found"));
        }
        Ok(Response::new(shop_opt.unwrap().into()))
    }
    async fn get_shop_by_shopkeeper_id(&self, _shopkeeper_id: Request<ShopkeeperId>) -> Result<Response<ShopInfo>, Status>{
        Err(Status::unimplemented("i don't think this is needed for now"))
    }
    async fn get_shop_by_shop_and_shopkeeper_id(&self, req: Request<ShopAndShopkeeperId>) -> Result<Response<ShopInfo>, Status>{
        println!("get_shop_by_shop_and_shopkeeper_id called");
        let ShopAndShopkeeperId{shop_id, shopkeeper_id} = req.into_inner();
        let shop_obj_id = ObjectId::parse_str(shop_id);
        if shop_obj_id.is_err() {
            return Err(Status::invalid_argument("invalid shop id"));
        }
        let shop_res = self.mongo_manager.shop_col.find_one(
            doc!{
                "_id": shop_obj_id.unwrap(),
                "user_id": shopkeeper_id
            }, None
        ).await;
        if shop_res.is_err() {
            return Err(Status::internal("unable to access shop db"));
        }
        let shop_opt = shop_res.unwrap();
        if shop_opt.is_none() {
            return Err(Status::not_found("shop not found"));
        }
        Ok(Response::new(shop_opt.unwrap().into()))
    }
}






