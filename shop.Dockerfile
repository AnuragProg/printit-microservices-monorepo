FROM rust:1.76

# setting working directory
WORKDIR /usr/app/shop

# install protobuf compiler for generation
RUN apt-get update
RUN apt-get install -y protobuf-compiler

# copying source code
COPY ./shop .

# copying proto definitions
COPY ./proto_def ../proto_def

# install corresponding dependencies
RUN cargo build #--release # commented for development purposes

CMD ["./target/debug/shop"]
