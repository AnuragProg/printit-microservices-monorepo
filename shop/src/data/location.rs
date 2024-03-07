use serde::{Serialize, Deserialize};



#[derive(Debug, Serialize, Deserialize)]
pub struct Location{
    #[serde(rename="type")]
    pub location_type: String,
    pub coordinates: [f64;2] // [lng, lat]
}
