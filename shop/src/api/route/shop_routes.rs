use rocket::State;
use rocket::response::{content::RawJson, status::Custom};
use rocket::{get, post, patch};
use rocket::http::Status;
use rocket::serde::json::{json, Json};
use crate::api::guards::auth_guard::{AuthUser, AuthShopkeeper};
use crate::client::mongo_client::MongoManager;
use crate::data::shop::{Shop, ShopBody, ShopUpdateBody};
use crate::data::location::Location;
use mongodb::{options, bson::{doc, oid::ObjectId}};


// TODO for the whole service
// add redis caching layer to cache the details and update them properly timely as updates comes in


#[get("/<shop_id>")]
pub async fn get_shop_details(
    _auth_user: AuthUser,
    mongo_manager: &State<MongoManager>,
    shop_id: &str
) -> Custom<RawJson<String>> {

    // verify shop_id correctness
    let id = ObjectId::parse_str(shop_id);
    if id.is_err() {
        return Custom(Status::BadRequest, RawJson(json!({"message": "invalid shop id"}).to_string()));
    }

    // retrieve shop details
    let shop_result = mongo_manager.shop_col.find_one(doc!{"_id": id.unwrap()}, None).await;
    if let Err(err) = shop_result{
        return Custom(Status::InternalServerError, RawJson(json!({"message": err.to_string()}).to_string()));
    }
    let shop_opt = shop_result.unwrap();
    if shop_opt.is_none() {
        return Custom(Status::NotFound, RawJson(json!({"message": "shop not found"}).to_string()));
    }
    Custom(Status::Ok, RawJson(json!({"shop": shop_opt.unwrap()}).to_string()))
}


// TODO add one shop guard as well
#[post("/", data="<shop_details>")]
pub async fn create_shop(
    auth_shopkeeper: AuthShopkeeper,
    mongo_manager: &State<MongoManager>,
    shop_details: Json<ShopBody>
) -> Custom<RawJson<String>>{

    // extract shop details
    let ShopBody{ name, contact, email, location } = shop_details.into_inner();

    // create shop document
    let now = bson::DateTime::from_chrono(chrono::Utc::now());
    let shop = Shop{
        _id: ObjectId::new(),
        user_id: auth_shopkeeper.user.id,
        name,
        contact,
        email,
        location: Location {
            location_type: String::from("Point"),
            coordinates: [location[0], location[1]]
        },
        created_at: now,
        updated_at: now,
    };

    // insert shop
    let result = mongo_manager.shop_col.insert_one(&shop, None).await;
    if let Err(err) = result{
        return Custom(Status::InternalServerError, RawJson(json!({"message": err.to_string()}).to_string()));
    }
    Custom(Status::Ok, RawJson(json!({"message": "shop created successfully", "shop": shop}).to_string()))
}


#[patch("/<shop_id>", data="<shop_update_details>")]
pub async fn update_shop(
    auth_shopkeeper: AuthShopkeeper,
    mongo_manager: &State<MongoManager>,
    shop_id: &str,
    shop_update_details: Json<ShopUpdateBody>
) -> Custom<RawJson<String>>{

    let user_id = auth_shopkeeper.user.id;

    // verify shop_id correctness
    let shop_obj_id_res = ObjectId::parse_str(shop_id);
    if shop_obj_id_res.is_err() {
        return Custom(Status::BadRequest, RawJson(json!({"message": "invalid shop id"}).to_string()));
    }
    let shop_obj_id = shop_obj_id_res.unwrap();

    // build shop query filter (we are adding user id to make sure only user's shop get's updated)
    let shop_query = doc!{"_id": shop_obj_id, "user_id": user_id};

    // check if shop exists
    let shop_exists_res = mongo_manager.shop_col.find_one(shop_query.clone(), None).await;
    if let Err(err) = shop_exists_res{
        return Custom(Status::InternalServerError, RawJson(json!({"message": err.to_string()}).to_string()));
    }
    let shop_exists = shop_exists_res.unwrap();
    if shop_exists.is_none() {
        return Custom(Status::NotFound, RawJson(json!({"message": "shop not found"}).to_string()));
    }

    // build the update doc
    let mut update_doc = doc!{"updated_at": bson::DateTime::from_chrono(chrono::Utc::now())};
    let update_details = shop_update_details.into_inner();
    if let Some(name) = update_details.name {
        update_doc.insert("name", name);
    }
    if let Some(contact) = update_details.contact {
        update_doc.insert("contact", contact);
    }
    if let Some(email) = update_details.email {
        update_doc.insert("email", email);
    }
    if let Some(location) = update_details.location {
        update_doc.insert("location", doc!{
            "type": "Point",
            "coordinates": vec![location[0], location[1]]
        });
    }
    println!("doc = {}", update_doc);

    // update the shop details
    let update_result = mongo_manager.shop_col.find_one_and_update(
        shop_query,
        doc!{"$set": update_doc},
        options::FindOneAndUpdateOptions::builder().return_document(options::ReturnDocument::After).build()
    ).await;
    if let Err(err) = update_result{
        return Custom(Status::InternalServerError, RawJson(json!({"message": err.to_string()}).to_string()));
    }
    let shop = update_result.unwrap();
    if shop.is_none() {
        return Custom(Status::NotFound, RawJson(json!({"message": "shop not found"}).to_string()));
    }
    Custom(Status::Ok, RawJson(json!({"message": "shop data updated successfully", "shop": shop.unwrap()}).to_string()))
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

