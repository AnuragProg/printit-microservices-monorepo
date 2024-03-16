use serde::{Serialize, Deserialize};
use crate::service::shop::{Location as LocationProto};
use std::convert::From;



#[derive(Debug, Serialize, Deserialize)]
pub struct Location{
    #[serde(rename="type")]
    pub location_type: String,
    pub coordinates: [f64;2] // [lng, lat]
}

impl From<Location> for LocationProto{
    fn from(location: Location) -> Self {
        LocationProto {
            lng: location.coordinates[0],
            lat: location.coordinates[1]
        }
    }
}
