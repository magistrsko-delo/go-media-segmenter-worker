package Worker

import (
	"encoding/json"
	"fmt"
	"log"
	"main/HttpUtils"
	"main/Models"
	"main/ffmpeg"
	"main/grpc_client"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Worker struct {
	RabbitMQ *RabbitMqConnection
	mediaMetadataGrpcClient *grpc_client.MediaMetadataClient
	mediaChunksClient *grpc_client.MediaChunksClient
	env *Models.Env
}


func (worker *Worker) Work()  {
	forever := make(chan bool)
	go func() {
		for d := range worker.RabbitMQ.msgs {
			log.Printf("Received a message: %s", d.Body)

			mediaMetadata := &Models.MediaMetadata{}
			err := json.Unmarshal([]byte(d.Body), mediaMetadata)
			if err != nil{
				log.Println(err)
			}

			fileUrl := worker.env.MediaManagerUrl + "v1/mediaManager/" + mediaMetadata.AwsBucketWholeMedia + "/" + mediaMetadata.AwsStorageNameWholeMedia
			err = HttpUtils.DownloadFile("./assets/" + mediaMetadata.AwsStorageNameWholeMedia, fileUrl)
			if err != nil {
				log.Println(err)
			}
			// exec ffmpeg commands
			ffmpeg := &ffmpeg.FFmpeg{}
				// get file frame rate
			frameRate, err := ffmpeg.ExecFFprobeCommand([]string{"-v", "0", "-of", "csv=p=0", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", "./assets/"+mediaMetadata.AwsStorageNameWholeMedia})
			resolution := "1920x1080"
			if err != nil {
				log.Println(err)
			}
			frames, _ := strconv.Atoi(strings.Split(frameRate, "/")[0])
			fmt.Println("FRAME RATE: ", frames)

			log.Println("CREATING VIDEO SEGMENTS")
			cmdArgs := []string{"-i", "./assets/"+ mediaMetadata.AwsStorageNameWholeMedia, "-vf", "scale=w=1920:h=1080",
				"-c:a", "aac", "-ar", "48000", "-b:a", "128k", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-g", strconv.Itoa(frames * 5), "-keyint_min", strconv.Itoa(frames * 5),
				"-sc_threshold", "0", "-b:v", "5000k", "-maxrate", "5350k", "-bufsize", "7500", "-b:a", "192k", "-hls_segment_filename", "./assets/chunks/1080p_%03d.ts", "./assets/chunks/1080p.m3u8"}

			err = ffmpeg.ExecFFmpegCommand(cmdArgs)

			if err != nil {
				log.Println(err)
			}
			var files []string
			root := "./assets/chunks"
			err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				files = append(files, path)
				return err
			})
			if err != nil {
				log.Println(err)
			}

			position := 0
			for index, file := range files {
				if index == 0 || strings.Contains(file, ".gitkeep") {
					continue
				}

				if strings.Contains(file, ".m3u8") {
					worker.removeFile(file)
					continue
				}

				fmt.Println(file)
				filePathArray := []string{}
				if worker.env.Env == "live" {
					filePathArray = strings.Split(file, "/")
				} else {
					filePathArray = strings.Split(file, "\\")
				}
				// (path string, mediaName string, awsBucket string, awsStorageHost string
				err = HttpUtils.UploadFile(file, filePathArray[len(filePathArray) - 1], mediaMetadata.AwsBucketWholeMedia, worker.env.AawsStorageUrl)
				if err != nil {
					log.Println(err)
				}

				length, err := ffmpeg.ExecFFprobeCommand([]string{"-i", file, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0"})
				if err != nil {
					log.Println(err)
				}

				lengthFloat, _ :=  strconv.ParseFloat(worker.standardizeSpaces(length), 64)
				lengthFloat = math.Ceil(lengthFloat*1000000)/1000000
				fmt.Println("length: ", lengthFloat)
				chunkMetadata := Models.NewChunkMetadata(mediaMetadata.AwsBucketWholeMedia, filePathArray[len(filePathArray) - 1], lengthFloat, mediaMetadata.MediaId, resolution, position)
				_, err = worker.mediaChunksClient.UpdateMediaMetadata(chunkMetadata)
				if err != nil {
					log.Println(err)
				}

				position++
				worker.removeFile(file)
			}

			worker.removeFile("./assets/" + mediaMetadata.AwsStorageNameWholeMedia)

			_, err = worker.mediaMetadataGrpcClient.UpdateMediaMetadata(mediaMetadata)
			if err != nil {
				log.Println(err)
			}

			log.Printf("Done")
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (worker *Worker) removeFile(path string)  {
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}

func (worker *Worker) standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func InitWorker() *Worker  {

	return &Worker{
		RabbitMQ: 					initRabbitMqConnection(Models.GetEnvStruct()),
		mediaMetadataGrpcClient: 	grpc_client.InitMediaMetadataGrpcClient(),
		mediaChunksClient:			grpc_client.InitChunkMetadataClient(),
		env:      					Models.GetEnvStruct(),
	}
}

