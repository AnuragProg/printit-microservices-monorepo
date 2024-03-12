use serde::{Serialize, Deserialize};
use crate::data::location::Location;
use mongodb::bson::oid::ObjectId;
use bson::DateTime;
use crate::services::shop_service::{ShopInfo as ShopInfoProto, Location as LocationProto};

#[derive(Debug, Serialize, Deserialize)]
pub struct Shop{
    pub _id: ObjectId,
    pub user_id: String,
    pub name: String,
    pub contact: String,
    pub email: String,
    pub location: Location,
    pub created_at: DateTime,
    pub updated_at: DateTime,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ShopResponse{
    pub _id: String,
    pub user_id: String,
    pub name: String,
    pub contact: String,
    pub email: String,
    pub location: Location,
    pub created_at: String,
    pub updated_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ShopBody{
    pub name: String,
    pub contact: String,
    pub email: String,
    pub location: [f64;2]
}


#[derive(Debug, Serialize, Deserialize)]
pub struct ShopUpdateBody{
    pub name: Option<String>,
    pub contact: Option<String>,
    pub email: Option<String>,
    pub location: Option<[f64;2]>
}





impl From<Shop> for ShopInfoProto{
    fn from(val: Shop) -> Self {
        ShopInfoProto{
            id: val._id.to_hex(),
            user_id: val.user_id,
            name: val.name,
            contact: val.contact,
            email: val.email,
            location: Some(LocationProto::from(val.location)),
            created_at: val.created_at.try_to_rfc3339_string().unwrap(),
            updated_at: val.updated_at.try_to_rfc3339_string().unwrap(),
        }
    }
}

impl From<Shop> for ShopResponse{
    fn from(val: Shop) -> Self {
        ShopResponse{
            _id: val._id.to_hex(),
            user_id: val.user_id,
            name: val.name,
            contact: val.contact,
            email: val.email,
            location: val.location,
            created_at: val.created_at.try_to_rfc3339_string().unwrap(),
            updated_at: val.updated_at.try_to_rfc3339_string().unwrap(),
        }
    }
}



