package oss

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

var (
	errImageTooLarge = errors.New("image too large")
	errIsNotImage    = errors.New("not image")
)

type OssClient struct {
	bucket *oss.Bucket
	option *Option
}

// Option oss client option
func NewOssClient(option *Option) *OssClient {
	option.init()

	client, err := oss.New(option.Endpoint, option.AccessKeyId, option.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	bucket, err := client.Bucket(option.BucketName)
	if err != nil {
		panic(err)
	}
	return &OssClient{
		bucket: bucket,
		option: option,
	}
}

// 上传本地文件路径
func (o *OssClient) UploadFilePath(toPath, filePath string) error {
	return o.bucket.PutObjectFromFile(toPath, filePath)
}

// 上传文件流
func (o *OssClient) UploadFileStream(toPath string, stream io.Reader) error {
	return o.bucket.PutObject(toPath, stream)
}

// 上传图片
func (o *OssClient) UploadImage(filename, stream string) (string, error) {
	fileContentPosition := strings.Index(stream, ",")
	uploadBaseString := stream[fileContentPosition+1:]
	uploadString, err := base64.StdEncoding.DecodeString(uploadBaseString)
	if err != nil {
		return "", err
	}

	if len(uploadString) > o.option.MaxFileSize {
		return "", errImageTooLarge
	}
	extName, isImage := IsImage(filename)
	if !isImage {
		return "", errIsNotImage
	}

	filename = uuid.New().String() + extName
	err = o.bucket.PutObject(o.option.BasePath+filename, strings.NewReader(string(uploadString)))
	if err != nil {
		return "", err
	}

	return o.option.prefix + o.bucket.BucketName + "." + o.bucket.Client.Config.Endpoint + "/" + o.option.BasePath + filename, nil
}

// 获取文件类型
func GetFileContentType(out *os.File) (string, error) {
	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

// 是否为图片
func IsImage(filename string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return ext, true
	}
	return "", false
}
