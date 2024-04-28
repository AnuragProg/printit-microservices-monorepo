FROM node:21-slim

# installing required tooling
RUN npm i -g pnpm protoc-gen-ts
#RUN apk add --no-cache protobuf #(not working with slim image)
RUN apt update && apt install -y protobuf-compiler

WORKDIR /usr/app

# copy package and pnpm files needed for downloading packages
COPY ./notification/package.json ./
COPY ./notification/pnpm-lock.yaml ./

# install packages
RUN pnpm install

# copy rest of source code
COPY ./notification .

# make necessary directories for proto stubs
RUN mkdir -p ./src/proto_gen/authentication
RUN mkdir -p ./src/proto_gen/order
RUN mkdir -p ./src/proto_gen/shop

# copy proto definitions
COPY ./proto_def ../proto_def

# generate proto files
RUN protoc --ts_out=./src/proto_gen/authentication --proto_path=../proto_def/authentication auth.proto
RUN protoc --ts_out=./src/proto_gen/order --proto_path=../proto_def/order order.proto
RUN protoc --ts_out=./src/proto_gen/shop --proto_path=../proto_def/shop shop.proto

RUN pnpm run build

CMD ["pnpm", "start"]
