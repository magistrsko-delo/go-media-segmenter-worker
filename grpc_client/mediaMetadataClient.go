package grpc_client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	pbMediaMetadata "main/proto/media_metadata"
	)


type MediaMetadataClient struct {
	conn *grpc.ClientConn
	client pbMediaMetadata.MediaMetadataClient
}

func (mediaMetadata *MediaMetadataClient) CreateNewMediaMetadata()  {
	response, err := mediaMetadata.client.NewMediaMetadata(context.Background(), &pbMediaMetadata.CreateNewMediaMetadataRequest{
		Name:                     "test_name",
		SiteName:                 "test_site",
		Length:                   1,
		Status:                   0,
		Thumbnail:                "waw",
		ProjectId:                -1,
		AwsBucketWholeMedia:      "aws bucket",
		AwsStorageNameWholeMedia: "whole_media",
	})

	if err != nil {
		log.Fatalf("could create new  media metadata: %v", err)
	}

	fmt.Println("RESPONSE")
	log.Println(response.GetMediaId())
	log.Println(response.GetCreatedAt())
	log.Println(response.GetAwsBucketWholeMedia())
}


func InitMediaMetadataGrpcClient() *MediaMetadataClient  {
	fmt.Println("CONNECTING")
	conn, err := grpc.Dial("localhost:9002", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION")

	client := pbMediaMetadata.NewMediaMetadataClient(conn)
	return &MediaMetadataClient{
		conn:    conn,
		client:  client,
	}

}
