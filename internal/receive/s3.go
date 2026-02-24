package receive

import (
	"log"
	"netinfo/internal/preload"
	"netinfo/internal/wireguard"
	"time"

	"github.com/unix755/xtools/xS3"
)

// getNetInfoFromS3 从 s3 服务器获取指定 id 的网络信息
func getNetInfoFromS3(endpoint string, region string, accessKeyId string, secretAccessKey string, stsToken string, pathStyle bool, allowInsecure bool, bucket string, objectPath string, encryptionKey []byte) (p *preload.Preload, err error) {
	c := xS3.NewS3Client(endpoint, region, accessKeyId, secretAccessKey, stsToken, pathStyle, allowInsecure)
	d, err := c.GetObject(bucket, objectPath)
	if err != nil {
		return nil, err
	}

	// 读取从 s3 服务器下载的数据流
	return preload.Unmarshal(d, "json", encryptionKey)
}

func FromS3(endpoint string, region string, accessKeyId string, secretAccessKey string, stsToken string, pathStyle bool, allowInsecure bool, bucket string, objectPath string, encryptionKey []byte, remoteInterface string, wgInterface string, wgPeerKey string) (err error) {
	p, err := getNetInfoFromS3(endpoint, region, accessKeyId, secretAccessKey, stsToken, pathStyle, allowInsecure, bucket, objectPath, encryptionKey)
	if err != nil {
		return err
	}
	publicIP, err := p.GetPublicIP(remoteInterface)
	if err != nil {
		return err
	}
	return wireguard.UpdateEndpoint(wgInterface, wgPeerKey, publicIP, -1)
}

func FromS3Loop(endpoint string, region string, accessKeyId string, secretAccessKey string, stsToken string, pathStyle bool, allowInsecure bool, bucket string, objectPath string, encryptionKey []byte, remoteInterface string, wgInterface string, wgPeerKey string, interval time.Duration) {
	for {
		err := FromS3(endpoint, region, accessKeyId, secretAccessKey, stsToken, pathStyle, allowInsecure, bucket, objectPath, encryptionKey, remoteInterface, wgInterface, wgPeerKey)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(interval)
	}
}
