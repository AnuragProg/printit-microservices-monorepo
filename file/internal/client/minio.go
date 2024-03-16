package client

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	consts "github.com/AnuragProg/printit-microservices-monorepo/internal/constant"
)



func GetMinioClient(minioURI, minioServerAccessKey, minioServerSecretKey string) (*minio.Client, error){

	// create client
	minioClient, err := minio.New(minioURI, &minio.Options{
		Creds: credentials.NewStaticV4(minioServerAccessKey, minioServerSecretKey, ""),
		Transport: &http.Transport{
			MaxIdleConns: 100,
			IdleConnTimeout: 60*time.Second,
		},
	})
	if err != nil{
		return nil, err
	}

	// check minio connection
	if _, err := minioClient.ListBuckets(context.Background()); err != nil{
		return nil, err
	}

	// Create required buckets
	minioClient.MakeBucket(context.Background(), consts.FILE_BUCKET, minio.MakeBucketOptions{})
	log.Println("Connected to Minio...")

	return minioClient, nil
}
