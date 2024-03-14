

fn main() -> Result<(), Box<dyn std::error::Error>>{
    tonic_build::compile_protos("../proto_def/authentication/auth.proto")?;
    tonic_build::compile_protos("../proto_def/shop/shop.proto")?;
    tonic_build::compile_protos("../proto_def/price/price.proto")?;
    Ok(())
}

