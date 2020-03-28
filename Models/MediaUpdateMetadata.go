package Models

import "strings"

type MediaUpdateMetadata struct {
	MediaId int `json:"mediaId"`
	Name string `json:"name"`
	SiteName string `json:"siteName"`
	Length int `json:"length"`
	Status int `json:"status"`
	Thumbnail string `json:"thumbnail"`
	ProjectId int `json:"projectId"`
	AwsBucketWholeMedia string `json:"awsBucketWholeMedia"`
	AwsStorageNameWholeMedia string `json:"awsStorageNameWholeMedia"`
	Keywords string `json:"keywords"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func InitMediaUpdateMetadata(mediaMetadata *MediaMetadata) *MediaUpdateMetadata  {
	return &MediaUpdateMetadata{
		MediaId:                  mediaMetadata.MediaId,
		Name:                     mediaMetadata.Name,
		SiteName:                 mediaMetadata.SiteName,
		Length:                   mediaMetadata.Length,
		Status:                   3,
		Thumbnail:                mediaMetadata.Thumbnail,
		ProjectId:                mediaMetadata.ProjectId,
		AwsBucketWholeMedia:      mediaMetadata.AwsBucketWholeMedia,
		AwsStorageNameWholeMedia: mediaMetadata.AwsStorageNameWholeMedia,
		Keywords:                 strings.Join(mediaMetadata.Keywords, ","),
		CreatedAt:                mediaMetadata.CreatedAt,
		UpdatedAt:                mediaMetadata.UpdatedAt,
	}
}
