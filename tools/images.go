package tools

import (
	"bytes"
	"context"
	"encoding/base64"
	"ghost-images/api"
	"io"
	"os"
	"path/filepath"
	"time"
)

var uploadUrl = "/images/upload/"

func saveBase64Image(image string) (string, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return "", err
	}

	imageName := time.Now().Format("20060102150405") + ".png"
	imagePath := filepath.Join("images", imageName)

	err = os.MkdirAll("images", 0755)
	if err != nil {
		return "", err
	}

	file, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(imageBytes))
	if err != nil {
		return "", err
	}

	return imagePath, nil
}

// curl -X POST -F 'file=@/path/to/images/my-image.jpg' -F 'ref=path/to/images/my-image.jpg' -H "Authorization: 'Ghost $token'" -H "Accept-Version: $version" https://{admin_domain}/ghost/api/admin/images/upload/

func UploadImageBase64Image(base64Image string) (map[string]interface{}, error) {
	imagePath, err := saveBase64Image(base64Image)
	if err != nil {
		return nil, err
	}

	response, err := api.MakeGhostMultipartRequestCurl(uploadUrl, map[string]string{}, context.Background(), false, "POST", imagePath)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func UploadImageLocalPath(localPath string) (map[string]interface{}, error) {
	response, err := api.MakeGhostMultipartRequestCurl(uploadUrl, map[string]string{}, context.Background(), false, "POST", localPath)
	if err != nil {
		return nil, err
	}

	return response, nil
}
