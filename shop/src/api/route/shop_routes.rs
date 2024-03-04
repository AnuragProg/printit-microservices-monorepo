use rocket::{get, post, patch};


#[get("/<shop_id>")]
async fn get_shop_details(shop_id: &str){
}


#[post("/")]
async fn create_shop(){
}


#[patch("/<shop_id>")]
async fn update_shop(shop_id: &str){
}


#[get("/live-location")]
async fn shop_location_ws(ws: ws::WebSocket) -> ws::Stream!['static]{
    ws::Stream! { ws =>
        for await message in ws {
            yield message?;
        }
    }
}

