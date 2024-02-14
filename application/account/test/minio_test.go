package test

import (
	"log"
	"testing"
	"time"
	"user/pkg/oss"
)

const (
	filePath         = "SampleVideo_1280x720_1mb.mp4"
	objectName       = "SampleVideo_1280x720_1mb.mp4"
	downloadFilePath = "SampleVideo_1280x720_1mb.mp4"
)

type MinioKeys struct {
	MinioEndPoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

// MinioInitTest 初始化minio client
func MinioInitTest(c MinioKeys) *oss.MinioClient {
	// 加载minio的endpoint等

	log.Println(c.MinioEndPoint, c.MinioAccessKey, c.MinioSecretKey)

	client, err := oss.NewMinioClient(c.MinioEndPoint, c.MinioAccessKey, c.MinioSecretKey)
	if err != nil {
		log.Println("init minio fail")
		return nil
	}
	return client
}

var minioKeys = MinioKeys{
	"192.168.2.184:9000",
	// "pupcFveqqthin0Q3eGsw",
	// "vJV6ryU3Qzkla6PeZta7MrfzpLxfduItjYgyzCns",
	"OzEzhBcDDbgsaFirC0td",
	"nQPMazFptbNzpl8NLeRFtkz2LkR4AmOzMOhwqEqu",
	"tiktok",
}

// 测试上传文件
func TestMinioUpload(t *testing.T) {
	client := MinioInitTest(minioKeys)
	err := client.UploadFile(minioKeys.MinioBucketName, objectName, filePath)
	if err != nil {
		t.Error("upload file: ", err)
	}
}

// 测试下载文件
func TestMinioDownloadFile(t *testing.T) {
	client := MinioInitTest(minioKeys)
	err := client.DownloadFile(minioKeys.MinioBucketName, objectName, downloadFilePath)
	if err != nil {
		t.Error("download file: ", err)
	}
}

// 测试列出bucket下所有的对象
func TestListObjects(t *testing.T) {
	client := MinioInitTest(minioKeys)
	objects, err := client.ListObjects(minioKeys.MinioBucketName, "")
	if err != nil {
		t.Error("list object: ", err)
	}
	log.Println("objects: ", objects)
}

// 删除对象
func TestDeleteFile(t *testing.T) {
	client := MinioInitTest(minioKeys)
	ret, err := client.DeleteFile(minioKeys.MinioBucketName, objectName)
	if err != nil {
		t.Error("delete object: ", ret, err)
	}
	log.Println("delete object: ", ret)
}

// 获取对象，返回url,
func TestGetPresignedGetObject(t *testing.T) {
	client := MinioInitTest(minioKeys)
	object, err := client.GetPresignedGetObject(minioKeys.MinioBucketName, objectName, 24*time.Hour)
	if err != nil {
		t.Error("GetPresignedGetObject: ", err)
	}
	log.Println("GetPresignedGetObject: ", object)
}
