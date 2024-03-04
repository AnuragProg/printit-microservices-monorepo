use rocket::fairing::{Fairing, Info, Kind};
use rocket::{Request, Data, Response};
use rocket::http::{Method, Status, ContentType};
use crate::client::auth_grpc_client::AuthGrpcManager;
use std::sync::Arc;
use regex::Regex;

pub struct ShopkeeperAuthMiddleware{
    auth_grpc_manager: Arc<AuthGrpcManager>,
    routes: Vec<Regex>
}

impl ShopkeeperAuthMiddleware{
    pub fn new(auth_grpc_manager: Arc<AuthGrpcManager>, routes: Vec<String>) -> Self{
        ShopkeeperAuthMiddleware{
            auth_grpc_manager,
            routes: routes.into_iter().map(|route| Regex::new(route.as_str()).unwrap()).collect()
        }
    }
}



// dump this code and use request guards to implement authentication
#[rocket::async_trait]
impl Fairing for ShopkeeperAuthMiddleware{
    fn info(&self) -> Info{
        Info{
            name: "Shopkeeper auth middlware",
            kind: Kind::Request
        }
    }

    async fn on_request(&self, request: &mut Request<'_>, _: &mut Data<'_>) {
        if self.routes.iter().any(|route| route.is_match(request.uri().path().as_str())) {
            if let Some(token) = request.headers().get_one("authentication") {

            }else{
                let message = "Unauthorized";
                let response = Response::build()
                    .status(Status::Unauthorized)
                    .header(ContentType::Plain)
                    .sized_body(message.len(), std::io::Cursor::new(message))
                    .finalize();
            }
        }
    }
}
