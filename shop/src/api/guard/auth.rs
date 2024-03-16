use rocket::State;
use rocket::http::Status;
use rocket::request::{FromRequest, Request, Outcome};
use crate::client::auth_grpc_client::AuthGrpcManager;
use crate::client::auth_grpc_client::{User, Token, UserType};

pub struct AuthUser{
    pub user: User
}

#[rocket::async_trait]
impl<'a> FromRequest<'a> for AuthUser{
    type Error = ();
    async fn from_request(request: &'a Request<'_>) -> Outcome<Self, Self::Error>{
        let mut auth_grpc_client = request.guard::<&State<AuthGrpcManager>>().await.unwrap().get_client();
        if let Some(auth_header) = request.headers().get_one("authorization"){
            let token = auth_header.trim_start_matches("Bearer ").to_string();
            if let Ok(user) = auth_grpc_client.verify_token(Token{token}).await {
                return Outcome::Success(AuthUser{user: user.into_inner()});
            }
        }
        Outcome::Error((Status::Unauthorized, ()))
    }
}

pub struct AuthShopkeeper{
    pub user: User
}

#[rocket::async_trait]
impl<'a> FromRequest<'a> for AuthShopkeeper{
    type Error = ();
    async fn from_request(request: &'a Request<'_>) -> Outcome<Self, Self::Error>{
        let mut auth_grpc_client = request.guard::<&State<AuthGrpcManager>>().await.unwrap().get_client();
        if let Some(auth_header) = request.headers().get_one("authorization"){
            let token = auth_header.trim_start_matches("Bearer ").to_string();
            if let Ok(user) = auth_grpc_client.verify_token(Token{token}).await{
                let user_info = user.into_inner();
                if user_info.user_type() == UserType::Shopkeeper{
                    return Outcome::Success(AuthShopkeeper{user: user_info});
                }
                println!("not a shopkeeper = {}", user_info.id);
            }
        }
        Outcome::Error((Status::Unauthorized, ()))
    }
}
