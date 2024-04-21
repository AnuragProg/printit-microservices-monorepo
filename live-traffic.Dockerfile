FROM node:21-slim

# installing required tooling
RUN npm i -g pnpm protoc-gen-ts
#RUN apk add --no-cache protobuf #(not working with slim image)
RUN apt update && apt install -y protobuf-compiler

WORKDIR /usr/app

# copy package and pnpm files needed for downloading packages
COPY ./live-traffic/package.json ./
COPY ./live-traffic/pnpm-lock.yaml ./

# install packages
RUN pnpm install

# copy rest of source code
COPY ./live-traffic .

# copy proto definitions
COPY ./proto_def ../proto_def

# generate proto files
RUN protoc --ts_out=./src/proto_gen/authentication --proto_path=../proto_def/authentication auth.proto
RUN protoc --ts_out=./src/proto_gen/order --proto_path=../proto_def/order order.proto

RUN pnpm run build

CMD ["pnpm", "start"]
