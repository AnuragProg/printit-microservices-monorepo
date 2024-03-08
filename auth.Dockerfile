FROM python:3.11-alpine

# setting working directory
WORKDIR /usr/app/authentication

# copy requirements file
COPY ./authentication/requirements.txt .

# install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# copying rest of the files
COPY ./authentication .

# copy proto definitions
COPY ./proto_def ../proto_def

# generate proto files
RUN mkdir -p ./src/proto_gen
RUN python -m grpc_tools.protoc --proto_path=../proto_def/authentication --python_out=./src/proto_gen --pyi_out=./src/proto_gen --grpc_python_out=./src/proto_gen auth.proto

# refreshing proto stubs and running the application
CMD python -u ./src/app.py

