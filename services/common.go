package services

import (
	"Twitta/global"
	"Twitta/pkg/constants"
	"Twitta/pkg/utils"
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/solost23/go_interface/gen_go/common"
	"github.com/solost23/go_interface/gen_go/oss"
	"io/ioutil"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"time"
)

func UploadImg(userId uint, folderName string, postFilename string, file *multipart.FileHeader) (string, error) {
	fileHandle, err := file.Open()
	if err != nil {
		return "", err
	}

	defer func() { _ = fileHandle.Close() }()

	b, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return "", err
	}

	return uploadImgOrVidBytes(userId, folderName, postFilename, b, "image")
}

func UploadVid(userId uint, folderName string, postFilename string, file *multipart.FileHeader) (string, error) {
	fileHandle, err := file.Open()
	if err != nil {
		return "", err
	}

	defer func() { _ = fileHandle.Close() }()

	b, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return "", err
	}

	return uploadImgOrVidBytes(userId, folderName, postFilename, b, "video")
}

func uploadImgOrVidBytes(userId uint, folderName string, postFileName string, fileBytes []byte, fileType string) (string, error) {
	if len(fileBytes) == 0 {
		return "", fmt.Errorf("upload image or video file is empty")
	}
	mime := strings.Split(mimetype.Detect(fileBytes).String(), " ")[0]
	if !strings.HasPrefix(mime, fileType) {
		return "", fmt.Errorf("invalid mime type: %s", mime)
	}

	filename := utils.NewMd5(
		time.Now().Format(constants.DateFormat)+
			fmt.Sprintf("%d", userId)+
			utils.NewMd5(string(fileBytes))+
			postFileName) + path.Ext(postFileName)
	url, err := upload(userId, folderName, filename, fileBytes)
	if err != nil {
		return "", err
	}

	return url, nil
}

func upload(userId uint, folder, filename string, fileBytes []byte) (url string, err error) {

	client := global.OSSSrvClient
	uploadResponse, err := client.Upload(context.Background(), &oss.UploadRequest{
		Header: &common.RequestHeader{
			Requester: strconv.Itoa(int(userId)),
			TraceId:   10000,
		},
		Folder:     folder,
		Key:        filename,
		Data:       fileBytes,
		UploadType: "static",
	})
	if err != nil {
		return "", err
	}
	return uploadResponse.Url, nil
}
