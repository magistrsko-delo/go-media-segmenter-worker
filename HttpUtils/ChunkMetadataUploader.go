package HttpUtils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"main/Models"
	"net/http"
)

func ChunkMetadataUpload(chunk *Models.ChunkMetadata, chunkMetadataHost string) error  {
	log.Println("CHUNK METADATA")
	log.Println(chunk)
	url := chunkMetadataHost + "v1/chunk/metadata/newMediaChunk"
	fmt.Println("URL chunk:>", url)

	metadataJsonString, err := json.Marshal(chunk)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(metadataJsonString))
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
