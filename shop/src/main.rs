mod api;
mod data;
mod services;
mod utils;
mod client;
use rocket::routes;
use std::sync::Arc;
use std::net::Ipv4Addr;
use tonic::transport::Server;
use client::mongo_client::MongoManager;
use api::route::shop_routes::{get_shop_details, create_shop, update_shop, shop_location_ws};
use client::auth_grpc_client::{authentication_client::AuthenticationClient, AuthGrpcManager, Empty};
use services::shop_service::{ShopService, shop_server::ShopServer};

#[rocket::main]
async fn main() -> Result<(), rocket::Error>{

    // consts
    let REST_PORT = std::env::var("REST_PORT").unwrap_or("3002".to_string());
    let GRPC_PORT = std::env::var("GRPC_PORT").unwrap_or("50053".to_string());
    let AUTH_GRPC_URI = std::env::var("AUTH_GRPC_URI").unwrap_or("http://[::1]:50051".to_string());
    let MONGO_URI = std::env::var("MONGO_URI").unwrap_or("mongodb://localhost:27017/?maxPoolSize=100".to_string());

    // create mongo manager
    let mongo_manager = Arc::new(MongoManager::new(MONGO_URI).await);

    // start grpc server
    let grpc_app = Server::builder()
        .add_service(ShopServer::new(ShopService::new(mongo_manager.clone())))
        .serve(format!("[::]:{GRPC_PORT}").parse().unwrap());

    // connecting to grpc clients
    let mut auth_grpc_client = AuthenticationClient::connect(AUTH_GRPC_URI).await.unwrap();
    auth_grpc_client.health_check(Empty{}).await.unwrap(); // check grpc client


    // creating rest config
    let rest_config = rocket::Config{
        address: Ipv4Addr::new(0,0,0,0).into(),
        port: REST_PORT.parse().unwrap(),
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
