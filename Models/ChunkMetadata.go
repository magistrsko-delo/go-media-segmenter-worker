package Models

type ChunkMetadata struct {
	AwsBucketName string `json:"awsBucketName"`
	AwsStorageName string `json:"awsStorageName"`
	Length float64 `json:"length"`
	MediaId int `json:"mediaId"`
	Resolution string `json:"resolution"`
	Position int `json:"position"`
}

func NewChunkMetadata(bucketName string, storageName string, length float64, mediaId int, resolution string, position int) *ChunkMetadata  {
	return &ChunkMetadata{
		AwsBucketName:  bucketName,
		AwsStorageName: storageName,
		Length:         length,
		MediaId:        mediaId,
		Resolution:     resolution,
		Position:       position,
	}
}