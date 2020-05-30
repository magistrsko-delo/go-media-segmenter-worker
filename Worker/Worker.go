package Worker

import (
	"encoding/json"
	"fmt"
	"log"
	"main/Http"
	"main/Models"
	"main/ffmpeg"
	"main/grpc_client"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type WorkerInterface interface {
	Work()
	getMediaFrameRate(*Models.MediaMetadata) (int, error )
	getSegmentLength(string) (float64, error)
	createFullHDVideoSegments(*Models.MediaMetadata, int) error
	getFilesPathsInDirectory() ([]string, error)
	handleMediaChunks([]string, *Models.MediaMetadata, string)  error
	removeFile(string)
	standardizeSpaces(string) string
}

type Worker struct {
	RabbitMQ *RabbitMqConnection
	mediaMetadataGrpcClient *grpc_client.MediaMetadataClient
	mediaChunksClient *grpc_client.MediaChunksClient
	awsStorageClient *grpc_client.AwsStorageClient
	ffmpeg *ffmpeg.FFmpeg
	mediaDowLoader *Http.MediaDownloader
	env *Models.Env
}


func (worker *Worker) Work()  {
	forever := make(chan bool)
	go func() {
		defer worker.awsStorageClient.Conn.Close()
		defer worker.mediaChunksClient.Conn.Close()
		defer worker.mediaMetadataGrpcClient.Conn.Close()
		for d := range worker.RabbitMQ.msgs {
			log.Printf("Received a message: %s", d.Body)
			isError := false

			mediaMetadata := &Models.MediaMetadata{}
			err := json.Unmarshal(d.Body, mediaMetadata)
			if err != nil{
				log.Println(err)
				isError = true
			}

			fileUrl := worker.env.AwsStorageUrl + "v1/awsStorage/media/" + mediaMetadata.AwsBucketWholeMedia + "/" + mediaMetadata.AwsStorageNameWholeMedia
			err = worker.mediaDowLoader.DownloadFile("./assets/" + mediaMetadata.AwsStorageNameWholeMedia, fileUrl)
			if err != nil {
				log.Println(err)
				isError = true
			}

			frames, err := worker.getMediaFrameRate(mediaMetadata)
			if err != nil {
				isError = true
			}

			//////////////////////////////////////////////// 1080p
			resolution := "1920x1080"

			err = worker.createFullHDVideoSegments(mediaMetadata, frames)
			if err != nil {
				isError = true
			}

			files, err := worker.getFilesPathsInDirectory()
			if err != nil {
				isError = true
			}

			err = worker.handleMediaChunks(files, mediaMetadata, resolution)
			if err != nil {
				isError = true
			}
			//////////////////////////480p
			resolution = "842x480"
			err = worker.create480pVideoSegments(mediaMetadata, frames)
			files, err = worker.getFilesPathsInDirectory()
			if err != nil {
				isError = true
			}
			err = worker.handleMediaChunks(files, mediaMetadata, resolution)
			if err != nil {
				isError = true
			}

			/////////////////////////

			/// media image
			thumbnail, err := worker.getMediaScreenShot(mediaMetadata)
			if err != nil {
				log.Println(err)
			}
			///

			worker.removeFile("./assets/" + mediaMetadata.AwsStorageNameWholeMedia)

			if !isError {
				mediaMetadata.Thumbnail = thumbnail
				_, err = worker.mediaMetadataGrpcClient.UpdateMediaMetadata(mediaMetadata)
				if err != nil {
					log.Println(err)
					isError = true
				}

				if !isError {
					log.Printf("Done")
					_ = d.Ack(false)
				}
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (worker *Worker) getMediaScreenShot(mediaMetadata *Models.MediaMetadata) (string, error) {
	log.Println("CREATING MEDIA SCREENSHOT")
	imageName := strconv.Itoa(mediaMetadata.MediaId) + "-" + strconv.Itoa(rand.Intn(1000000000000)) + "-" + mediaMetadata.AwsBucketWholeMedia + ".jpg"
	err := worker.ffmpeg.ExecFFmpegCommand([]string{"-ss", "00:00:01", "-i", "./assets/" + mediaMetadata.AwsStorageNameWholeMedia, "-vframes", "1", "-g:v", "2", "./assets/" + imageName})
	if err != nil {
		log.Println(err)
		return "" , err
	}
	_, err = worker.awsStorageClient.UploadMedia( "./assets/" + imageName, "mag20-images", imageName)  // TODO for later add this to configuration maybe..
	if err != nil {
		log.Println(err)
		return "" , err
	}

	worker.removeFile("./assets/" + imageName)
	return "v1/mediaManager/mag20-images/" + imageName, nil
}

func (worker *Worker) getMediaFrameRate(mediaMetadata *Models.MediaMetadata) (int, error ) {
	frameRate, err := worker.ffmpeg.ExecFFprobeCommand([]string{"-v", "0", "-of", "csv=p=0", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", "./assets/" + mediaMetadata.AwsStorageNameWholeMedia})
	if err != nil {
		log.Println(err)
		return -1, err
	}
	frames, _ := strconv.Atoi(strings.Split(frameRate, "/")[0])
	fmt.Println("FRAME RATE: ", frames)
	return frames, nil
}

func (worker *Worker) getSegmentLength(file string) (float64, error)  {
	length, err := worker.ffmpeg.ExecFFprobeCommand([]string{"-i", file, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0"})
	if err != nil {
		log.Println(err)
		return -1, err
	}
	lengthFloat, _ :=  strconv.ParseFloat(worker.standardizeSpaces(length), 64)
	lengthFloat = math.Ceil(lengthFloat*1000000)/1000000
	return lengthFloat, nil
}

func (worker *Worker) createFullHDVideoSegments(mediaMetadata *Models.MediaMetadata, frames int) error  {
	log.Println("CREATING VIDEO SEGMENTS 1080p")
	cmdArgs := []string{"-i", "./assets/" + mediaMetadata.AwsStorageNameWholeMedia, "-vf", "scale=w=1920:h=1080:force_original_aspect_ratio=decrease",
		"-c:a", "aac", "-ar", "48000", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-g", strconv.Itoa(frames * 5), "-keyint_min", strconv.Itoa(frames * 5),
		"-sc_threshold", "0", "-b:v", "5000k", "-maxrate", "5350k", "-bufsize", "7500k", "-b:a", "192k", "-hls_segment_filename", "./assets/chunks/1080p_%03d.ts", "./assets/chunks/1080p.m3u8"}

	err := worker.ffmpeg.ExecFFmpegCommand(cmdArgs)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (worker *Worker) create480pVideoSegments(mediaMetadata *Models.MediaMetadata, frames int) error  {
	log.Println("CREATING VIDEO SEGMENTS 480p")
	cmdArgs := []string{"-i", "./assets/"+ mediaMetadata.AwsStorageNameWholeMedia, "-vf", "scale=w=842:h=480:force_original_aspect_ratio=decrease",
		"-c:a", "aac", "-ar", "48000", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-g", strconv.Itoa(frames * 5), "-keyint_min", strconv.Itoa(frames * 5),
		"-sc_threshold", "0", "-b:v", "1400k", "-maxrate", "1400k", "-bufsize", "2100k", "-b:a", "128k", "-hls_segment_filename", "./assets/chunks/480p_%03d.ts", "./assets/chunks/480p.m3u8"}

	err := worker.ffmpeg.ExecFFmpegCommand(cmdArgs)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}


func (worker *Worker) getFilesPathsInDirectory() ([]string, error) {
	var files []string
	root := "./assets/chunks"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return err
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return files, nil
}

func (worker *Worker) handleMediaChunks(files []string, mediaMetadata *Models.MediaMetadata, resolution string)  error {
	position := 0
	for index, file := range files {
		if index == 0 || strings.Contains(file, ".gitkeep") {
			continue
		}

		if strings.Contains(file, ".m3u8") {
			worker.removeFile(file)
			continue
		}

		log.Println(file)
		filePathArray := []string{}
		if worker.env.Env == "live" {
			filePathArray = strings.Split(file, "/")
		} else {
			filePathArray = strings.Split(file, "\\")
		}

		_, err := worker.awsStorageClient.UploadMedia(file, mediaMetadata.AwsBucketWholeMedia, filePathArray[len(filePathArray) - 1])
		if err != nil {
			log.Println(err)
			return err
		}

		lengthFloat, err := worker.getSegmentLength(file)
		if err != nil {
			return err
		}

		log.Println("length: ", lengthFloat)
		chunkMetadata := Models.NewChunkMetadata(mediaMetadata.AwsBucketWholeMedia, filePathArray[len(filePathArray) - 1], lengthFloat, mediaMetadata.MediaId, resolution, position)

		_, err = worker.mediaChunksClient.UpdateMediaMetadata(chunkMetadata)
		if err != nil {
			log.Println(err)
			return err
		}
		position++
		worker.removeFile(file)
	}
	return nil
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
		awsStorageClient:			grpc_client.InitAwsStorageGrpcClient(),
		ffmpeg: 					&ffmpeg.FFmpeg{},
		mediaDowLoader:				&Http.MediaDownloader{},
		env:      					Models.GetEnvStruct(),
	}
}

