use tonic::transport::Channel;
use std::sync::Mutex;


tonic::include_proto!("auth_grpc");


pub struct AuthGrpcManager{
    client: Mutex<authentication_client::AuthenticationClient<Channel>>
}

impl AuthGrpcManager{
    pub fn new(client: authentication_client::AuthenticationClient<Channel>) -> Self{
        AuthGrpcManager{
            client: Mutex::new(client)
        }
    }
    pub fn get_client(&self) -> authentication_client::AuthenticationClient<Channel>{
        self.client.lock().unwrap().clone()
    }
}
