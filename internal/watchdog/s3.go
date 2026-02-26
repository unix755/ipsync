package watchdog

import (
	"ipsync/internal/preload"

	"github.com/unix755/xtools/xS3"
)

// FromS3 从 s3 服务器获取指定 id 的网络信息
func FromS3(endpoint string, region string, accessKeyId string, secretAccessKey string, stsToken string, pathStyle bool, skipTLSVerify bool, bucket string, objectPath string, encryptionKey []byte) (p preload.Preload, err error) {
	c := xS3.NewS3Client(endpoint, region, accessKeyId, secretAccessKey, stsToken, pathStyle, skipTLSVerify)
	d, err := c.GetObject(bucket, objectPath)
	if err != nil {
		return p, err
	}

	// 读取从 s3 服务器下载的数据流
	return preload.Unmarshal(d, "json", encryptionKey)
}
