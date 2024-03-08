FROM golang:1.21-alpine

# setting working directory
WORKDIR /usr/app/file

# copy mod and sum files
COPY ./file/go.mod ./file/go.sum ./

# install mod files
RUN go mod download && go mod verify

# copy proto definitions
COPY ./proto_def ../proto_def

# copying rest of the files
COPY ./file .

# installing protoc
RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# generate proto files
RUN mkdir -p proto_gen/authentication
RUN mkdir -p proto_gen/file
RUN protoc --go_out=./proto_gen/authentication --go_opt=paths=source_relative --go-grpc_out=./proto_gen/authentication --go-grpc_opt=paths=source_relative --proto_path=../proto_def/authentication ../proto_def/authentication/auth.proto
RUN protoc --go_out=./proto_gen/file --go_opt=paths=source_relative --go-grpc_out=./proto_gen/file --go-grpc_opt=paths=source_relative --proto_path=../proto_def/file ../proto_def/file/file.proto

# Build executable
RUN go build -o app ./cmd/

# running the application
CMD ["./app"]
