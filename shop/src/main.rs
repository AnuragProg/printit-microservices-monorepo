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
use services::shop_service::{ShopService, shop_server::ShopServer};
use tonic::transport::Server;

#[rocket::main]
async fn main() -> Result<(), rocket::Error>{

    // start grpc server
    let grpc_app = Server::builder()
        .add_service(ShopServer::new(ShopService))
        .serve("[::]:50052".parse().unwrap());

    // create mongo manager
    let mongo_manager = MongoManager::new().await;

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
        .manage(mongo_manager)
        .manage(AuthGrpcManager::new(auth_grpc_client))
        .configure(rest_config)
        .mount("/shop", routes![get_shop_details, create_shop, update_shop, shop_location_ws])
        .launch();

    tokio::join!(rest_app, grpc_app);
    Ok(())
}
