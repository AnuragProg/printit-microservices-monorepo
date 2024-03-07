tonic::include_proto!("shop_grpc");
use shop_server::Shop;
use tonic::{Request, Response, Status};


#[derive(Default, Debug)]
pub struct ShopService;


#[tonic::async_trait]
impl Shop for ShopService{
    async fn health_check(&self, _req: Request<Empty>) ->  Result<Response<Empty>, Status>{
        Ok(Response::new(Empty{}))
    }
}
