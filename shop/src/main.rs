


#[rocket::main]
async fn main() -> Result<(), rocket::Error>{
    let rest_app = rocket::build()
        .launch();
    tokio::join!(rest_app);
    Ok(())
}
