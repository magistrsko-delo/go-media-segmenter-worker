package Models

type MediaMetadata struct {
	MediaId int `json:"mediaId"`
	Name string `json:"name"`
	SiteName string `json:"siteName"`
	Length int `json:"length"`
	Status int `json:"status"`
	Thumbnail string `json:"thumbnail"`
	ProjectId int `json:"projectId"`
	AwsBucketWholeMedia string `json:"awsBucketWholeMedia"`
	AwsStorageNameWholeMedia string `json:"awsStorageNameWholeMedia"`
	Keywords []string `json:"keywords"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
