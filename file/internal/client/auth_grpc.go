package client

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



func GetAuthGrpcConnAndClient(authGrpcURI string) (*grpc.ClientConn, error){
	authGrpcConn, err := grpc.Dial(authGrpcURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		return nil, err
	}
	log.Println("Connected to Auth GRPC server...")
	return authGrpcConn, nil
}
