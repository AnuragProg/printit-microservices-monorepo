use serde::{Serialize, Deserialize};
use crate::data::location::Location;
use mongodb::bson::oid::ObjectId;

#[derive(Debug, Serialize, Deserialize)]
pub struct Shop{
    pub _id: ObjectId,
    pub user_id: String,
    pub name: String,
    pub contact: String,
    pub email: String,
    pub location: Location,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub updated_at: chrono::DateTime<chrono::Utc>,
}


#[derive(Debug, Serialize, Deserialize)]
pub struct ShopBody{
    pub name: String,
    pub contact: String,
    pub email: String,
    pub location: Location
}
