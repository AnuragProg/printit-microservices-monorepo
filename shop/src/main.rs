mod api;
mod data;
mod services;
mod utils;
mod client;

use std::net::Ipv4Addr;
//use api::middleware::auth_middlware::ShopkeeperAuthMiddleware; TODO figure out how to use request
//guards to authenticate requests

#[rocket::main]
async fn main() -> Result<(), rocket::Error>{

    let rest_config = rocket::Config{
        address: Ipv4Addr::new(0,0,0,0).into(),
        port: 3002,
        ..rocket::Config::default()
    };
    let rest_app = rocket::build()
        .configure(rest_config)
        //.attach(ShopkeeperAuthMiddleware{})
        .launch();

    tokio::join!(rest_app);
    Ok(())
}
