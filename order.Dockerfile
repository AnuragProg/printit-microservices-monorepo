FROM golang:1.21-alpine

# setting working directory
WORKDIR /usr/app/order

# installing protoc
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# copy mod and sum files
COPY ./order/go.mod ./order/go.sum ./

# install mod files
RUN go mod download && go mod verify

# copy proto definitions
COPY ./proto_def ../proto_def

# copying rest of the files
COPY ./order .

# generate proto files
RUN mkdir -p proto_gen/authentication proto_gen/file proto_gen/shop proto_gen/price proto_gen/order
RUN protoc --go_out=./proto_gen/authentication --go_opt=paths=source_relative --go-grpc_out=./proto_gen/authentication --go-grpc_opt=paths=source_relative --proto_path=../proto_def/authentication ../proto_def/authentication/auth.proto
RUN protoc --go_out=./proto_gen/file --go_opt=paths=source_relative --go-grpc_out=./proto_gen/file --go-grpc_opt=paths=source_relative --proto_path=../proto_def/file ../proto_def/file/file.proto
RUN protoc --go_out=./proto_gen/shop --go_opt=paths=source_relative --go-grpc_out=./proto_gen/shop --go-grpc_opt=paths=source_relative --proto_path=../proto_def/shop ../proto_def/shop/shop.proto
RUN protoc --go_out=./proto_gen/price --go_opt=paths=source_relative --go-grpc_out=./proto_gen/price --go-grpc_opt=paths=source_relative --proto_path=../proto_def/price ../proto_def/price/price.proto
RUN protoc --go_out=./proto_gen/order --go_opt=paths=source_relative --go-grpc_out=./proto_gen/order --go-grpc_opt=paths=source_relative --proto_path=../proto_def/order ../proto_def/order/order.proto

# Build executable
RUN go build -o app ./cmd/

# running the application
CMD ["./app"]
