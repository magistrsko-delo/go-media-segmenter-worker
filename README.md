# go-media-segmenter-worker

Media chunks segmenter worker

## Linux build

```GOOS=linux GOARCH=amd64 go build .```



## Worker example message

```{"mediaId":4,"keywords":[],"createdAt":"1585063230614","updatedAt":"1585063230614","name":"video sample","siteName":"24ur","length":29,"status":0,"awsBucketWholeMedia":"video-sample1","awsStorageNameWholeMedia":"SampleVideo_1280x720_5mb.mp4"}```

```ffmpeg -i SampleVideo_1280x720_5mb.mp4 -vf scale=w=1280:h=720:force_original_aspect_ratio=decrease -c:a aac -ar 48000 -b:a 128k -c:v h264 -profile:v main -crf 20 -g 50 -keyint_min 50 -sc_threshold 0 -b:v 2500k -maxrate 2675k -bufsize 3750k -hls_segment_filename 720p_%03d.ts 720p.m3u8```





```
cmdArgs := []string{"-i", "./assets/"+mediaMetadata.AwsStorageNameWholeMedia, "-vf", "scale=w=1920:h=1080:force_original_aspect_ratio=decrease",
   "-c:a", "aac", "-ar", "48000", "-b:a", "128k", "-c:v", "h264", "-profile:v", "main", "-crf", "20", "-g", "50", "-keyint_min", "50",
   "-sc_threshold", "0", "-b:v", "5000k", "-maxrate", "5350k", "-bufsize", "7500", "-b:a", "192k", "-hls_segment_filename", "./assets/chunks/1080p_%03d.ts", "./assets/chunks/1080p.m3u8"}
```

##PROTOCOL BUFFER

```.env
protoc proto\helloworld.proto --go_out=plugins=grpc:.
```