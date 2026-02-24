package s3

import (
	"log"
	"netinfo/internal/send/preload"
	"time"

	"github.com/unix755/xtools/xS3"
)

func SendRequest(endpoint string, region string, accessKeyId string, secretAccessKey string, stsToken string, pathStyle bool, allowInsecure bool, bucket string, objectPath string, encryptionKey []byte) (location *string, err error) {
	// 获取负载
	p, err := preload.GetPreload(encryptionKey)
	if err != nil {
		return nil, err
	}

	// 使用 s3 协议上传负载
	c := xS3.NewS3Client(endpoint, region, accessKeyId, secretAccessKey, stsToken, pathStyle, allowInsecure)
	result, err := c.UploadObject(bucket, objectPath, p)
	if err != nil {
		return nil, err
	}
	return result.Location, nil
}

func SendRequestLoop(endpoint string, region string, accessKeyId string, secretAccessKey string, stsToken string, pathStyle bool, allowInsecure bool, bucket string, objectPath string, encryptionKey []byte, interval time.Duration) {
	for {
		location, err := SendRequest(endpoint, region, accessKeyId, secretAccessKey, stsToken, pathStyle, allowInsecure, bucket, objectPath, encryptionKey)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("upload to %s", *location)
		}
		time.Sleep(interval)
	}
}
