mod api;
mod data;
mod services;
mod utils;
mod client;
use rocket::routes;
use std::net::Ipv4Addr;
use api::route::shop_routes::{get_shop_details, create_shop, update_shop, shop_location_ws};
use client::auth_grpc_client::{authentication_client::AuthenticationClient, AuthGrpcManager, Empty};
use client::mongo_client::MongoManager;
use mongodb::Client;

#[rocket::main]
async fn main() -> Result<(), rocket::Error>{

    // connect to mongo
    let mongo_client = Client::with_uri_str("mongodb://localhost:27017/prinit?maxPoolSize=100").await.unwrap();

    // connecting to grpc clients
    let mut auth_grpc_client = AuthenticationClient::connect("http://[::1]:50051").await.unwrap();
    auth_grpc_client.health_check(Empty{}).await.unwrap(); // check grpc client


    // creating rest config
    let rest_config = rocket::Config{
        address: Ipv4Addr::new(0,0,0,0).into(),
        port: 3002,
        ..rocket::Config::default()
    };

    //starting rest app
    let rest_app = rocket::build()
        .manage(MongoManager::new(mongo_client))
        .manage(AuthGrpcManager::new(auth_grpc_client))
        .configure(rest_config)
        .mount("/shop", routes![get_shop_details, create_shop, update_shop, shop_location_ws])
        //.attach(ShopkeeperAuthMiddleware{})
        .launch();

    tokio::join!(rest_app);
    Ok(())
}
