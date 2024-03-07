use serde::{Serialize, Deserialize};
use crate::data::location::Location;
use mongodb::bson::oid::ObjectId;
use bson::DateTime;

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
