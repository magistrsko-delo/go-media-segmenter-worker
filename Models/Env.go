package Models

import (
	"fmt"
	"os"
)

var envStruct *Env

type Env struct {
	RabbitUser string
	RabbitPassword string
	RabbitQueue string
	RabbitHost string
	MediaManagerUrl string
	MediaMetadataUrl string
	AawsStorageUrl string
	ChunkMetadataUrl string
	Env string
	MediaMetadataGrpcServer string
	MediaMetadataGrpcPort string
	ChunkMetadataGrpcServer string
	ChunkMetadataGrpcPort string
}

func InitEnv()  {
	envStruct = &Env{
		RabbitUser:       			os.Getenv("RABBIT_USER"),
		RabbitPassword:   			os.Getenv("RABBIT_PASSWORD"),
		RabbitQueue:      			os.Getenv("RABBIT_QUEUE"),
		RabbitHost:       			os.Getenv("RABBIT_HOST"),
		MediaManagerUrl:  			os.Getenv("MEDIA_MANAGER_URL"),
		MediaMetadataUrl: 			os.Getenv("MEDIA_METADATA_URL"),
		AawsStorageUrl:   			os.Getenv("AWS_STORAGE_URL"),
		ChunkMetadataUrl: 			os.Getenv("CHUNK_METADATA_URL"),
		Env: 			  			os.Getenv("ENV"),
		MediaMetadataGrpcServer: 	os.Getenv("MEDIA_METADATA_GRPC_SERVER"),
		MediaMetadataGrpcPort:   	os.Getenv("MEDIA_METADATA_GRPC_PORT"),
		ChunkMetadataGrpcServer:  	os.Getenv("CHUNK_METADATA_GRPC_SERVER"),
		ChunkMetadataGrpcPort:		os.Getenv("CHUNK_METADATA_GRPC_PORT"),
	}
	fmt.Println(envStruct)
}

func GetEnvStruct() *Env  {
	return  envStruct
}