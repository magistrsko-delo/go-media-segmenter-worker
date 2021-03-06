package grpc_client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"main/Models"
	pbMediaChunks "main/proto/media_chunks"
)

type MediaChunksClient struct {
	Conn *grpc.ClientConn
	client pbMediaChunks.MediaMetadataClient
}

func (mediaChunksClient *MediaChunksClient) UpdateMediaMetadata(chunkMetadata *Models.ChunkMetadata) (*pbMediaChunks.MediaChunkInfoResponseRepeated, error)  {

	response, err := mediaChunksClient.client.NewMediaChunk(context.Background(), &pbMediaChunks.NewMediaChunkRequest{
		AwsBucketName:        chunkMetadata.AwsBucketName,
		AwsStorageName:       chunkMetadata.AwsStorageName,
		Length:               chunkMetadata.Length,
		MediaId:              int32(chunkMetadata.MediaId),
		Resolution:           chunkMetadata.Resolution,
		Position:             int32(chunkMetadata.Position),
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}


func InitChunkMetadataClient() *MediaChunksClient  {
	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING chunks metadata")

	conn, err := grpc.Dial(env.ChunkMetadataGrpcServer + ":" + env.ChunkMetadataGrpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION chunk metadata")

	client := pbMediaChunks.NewMediaMetadataClient(conn)
	return &MediaChunksClient{
		Conn:    conn,
		client:  client,
	}

}
