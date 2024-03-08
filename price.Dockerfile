FROM python:3.11-alpine

# setting working directory
WORKDIR /usr/app/price

# copy requirements file
COPY ./price/requirements.txt .

# install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# copying rest of the files
COPY ./price .

# copy proto definitions
COPY ./proto_def ../proto_def

# generate proto files
RUN mkdir -p ./src/proto_gen
RUN python -m grpc_tools.protoc --proto_path=../proto_def/authentication --python_out=./src/proto_gen --pyi_out=./src/proto_gen --grpc_python_out=./src/proto_gen auth.proto
RUN python -m grpc_tools.protoc --proto_path=../proto_def/shop --python_out=./src/proto_gen --pyi_out=./src/proto_gen --grpc_python_out=./src/proto_gen shop.proto
RUN python -m grpc_tools.protoc --proto_path=../proto_def/price --python_out=./src/proto_gen --pyi_out=./src/proto_gen --grpc_python_out=./src/proto_gen price.proto

# refreshing proto stubs and running the application
CMD python -u ./src/app.py

