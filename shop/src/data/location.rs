use serde::{Serialize, Deserialize};



#[derive(Debug, Serialize, Deserialize)]
pub struct Location{
    #[serde(rename="type")]
    location_type: String,
    coordinates: Vec<f64> // [lng, lat]
}
