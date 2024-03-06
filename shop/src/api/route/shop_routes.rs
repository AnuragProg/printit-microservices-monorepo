use rocket::State;
use rocket::response::{content::RawJson, status::Custom};
use rocket::{get, post, patch};
use rocket::http::Status;
use rocket::serde::json::json;
use rocket::serde::json::Json;
use crate::api::guards::auth_guard::{AuthUser, AuthShopkeeper};
use crate::client::mongo_client::MongoManager;
use crate::data::shop::{Shop, ShopBody};
use mongodb::bson::oid::ObjectId;


// TODO for the whole service
// add redis caching layer to cache the details and update them properly timely as updates comes in




#[get("/<shop_id>")]
pub async fn get_shop_details(
    auth_user: AuthUser,
    mongo_manager: &State<MongoManager>,
    shop_id: &str
) -> Custom<RawJson<String>> {
    todo!();
}


// TODO add one shop guard as well
#[post("/", data="<shop_details>")]
pub async fn create_shop(
    auth_shopkeeper: AuthShopkeeper,
    mongo_manager: &State<MongoManager>,
    shop_details: Json<ShopBody>
) -> Custom<RawJson<String>>{
    let now = chrono::Utc::now();
    let ShopBody{ name, contact, email, location } = shop_details.into_inner();
    let shop = Shop{
        _id: ObjectId::new(),
        user_id: auth_shopkeeper.user.id,
        name,
        contact,
        email,
        location,
        created_at: now,
        updated_at: now,
    };
    let result = mongo_manager.shop_col.insert_one(&shop, None).await;
    if let Err(err) = result{
        return Custom(Status::InternalServerError, RawJson(json!({"message": err.to_string()}).to_string()));
    }
    println!("inserted shop with id = {}", result.unwrap().inserted_id);
    Custom(Status::Ok, RawJson(json!({"message": "shop created successfully", "shop": shop}).to_string()))
}


#[patch("/<shop_id>")]
pub async fn update_shop(shop_id: &str) -> Custom<RawJson<String>>{
    todo!();
}


// TODO
// add logic to get nearby shops and return it as a json payload
#[get("/live-location")]
pub async fn shop_location_ws(ws: ws::WebSocket) -> ws::Stream!['static]{
    ws::Stream! { ws =>
        for await message in ws {
            yield message?;
        }
    }
}

