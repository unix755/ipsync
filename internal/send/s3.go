package send

import (
	"encoding/json"
	"ipsync/internal/cache"
	"ipsync/internal/preload"
	"log"
	"time"

	"github.com/unix755/xtools/xS3"
)

func ToS3(endpoint string, region string, accessKeyId string, secretAccessKey string, stsToken string, pathStyle bool, skipTLSVerify bool, bucket string, objectPath string, encryptionKey []byte) (location *string, err error) {
	// 获取负载
	p, err := preload.NewPreload()
	if err != nil {
		return nil, err
	}
	// 负载转换为加密比特流
	bytes, err := preload.Marshal(p, "json", encryptionKey)
	if err != nil {
		return nil, err
	}

	// 使用 s3 协议上传负载
	c := xS3.NewS3Client(endpoint, region, accessKeyId, secretAccessKey, stsToken, pathStyle, skipTLSVerify)
	result, err := c.UploadObject(bucket, objectPath, bytes)
	if err != nil {
		return nil, err
	}
	return result.Location, nil
}

func ToS3Loop(endpoint string, region string, accessKeyId string, secretAccessKey string, stsToken string, pathStyle bool, skipTLSVerify bool, bucket string, objectPath string, encryptionKey []byte, interval time.Duration) {
	for {
		// 获取 preload
		p, err := preload.NewPreload()
		if err != nil {
			log.Println(err)
		}
		// 获取网络界面
		bytes, err := json.Marshal(p.NetInterfaces)
		if err != nil {
			log.Println(err)
		}

		// 获取缓存中的 net_interfaces
		cacheNetInterfaces, _ := cache.Get("net_interfaces")

		if string(bytes) != cacheNetInterfaces {
			// 发送到文件
			location, err := ToS3(endpoint, region, accessKeyId, secretAccessKey, stsToken, pathStyle, skipTLSVerify, bucket, objectPath, encryptionKey)
			if err != nil {
				log.Println(err)
			} else {
				// 设置新的缓存 net_interfaces
				cache.Set("net_interfaces", string(bytes))
				log.Printf("new net interfaces found, upload to %s", *location)
			}
		} else {
			log.Println("new net interfaces not found, skip")
		}

		time.Sleep(interval)
	}
}
