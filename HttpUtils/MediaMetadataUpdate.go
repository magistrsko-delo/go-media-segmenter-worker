package HttpUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"main/Models"
	"net/http"
	"strconv"
)

func UpdateMediaMetadata(metadata *Models.MediaUpdateMetadata, mediaMetadataHost string) error  {
	log.Println("METADATA UPDATE")
	log.Println(metadata)
	url := mediaMetadataHost + "v1/media/metadata/" + strconv.Itoa(metadata.MediaId) + "/update"
	fmt.Println("URL:>", url)

	metadataJsonString, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(metadataJsonString))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	return nil
}