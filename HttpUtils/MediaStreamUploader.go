package HttpUtils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"main/Models"
	"mime/multipart"
	"net/http"
	"os"
)

func NewFileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func UploadFile(path string, mediaName string, awsBucket string, awsStorageHost string) error  {
	fmt.Println("uploading: " + awsStorageHost + "v1/awsStorage/media/" + awsBucket + "/" + mediaName)
	request, err := NewFileUploadRequest(awsStorageHost + "v1/awsStorage/media/" + awsBucket + "/" + mediaName, "mediaStream", path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer resp.Body.Close()
	uploadResponse := Models.UploadResponse{}
	_ = json.NewDecoder(resp.Body).Decode(&uploadResponse)
	if uploadResponse.Status != 200 || !uploadResponse.Data {
		return errors.New(string(uploadResponse.Status) + " : " + uploadResponse.Message)
	}
	return nil
}
