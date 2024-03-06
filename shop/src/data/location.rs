use serde::{Serialize, Deserialize};



#[derive(Debug, Serialize, Deserialize)]
pub struct Location{
    location_type: String,
    coordinates: Vec<f64>
}
