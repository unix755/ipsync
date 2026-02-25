package receive

import (
	"io"
	"ipsync/internal/preload"

	"github.com/unix755/xtools/xWebDAV"
)

// FromWebDAV 从 webdav 服务器获取指定 id 的网络信息
func FromWebDAV(endpoint string, username string, password string, skipTLSVerify bool, filepath string, encryptionKey []byte) (p preload.Preload, err error) {
	client, err := xWebDAV.NewClient(endpoint, username, password, skipTLSVerify)
	if err != nil {
		return p, err
	}
	response, err := client.Download(filepath)
	if err != nil {
		return p, err
	}

	// 读取从 webdav 服务器下载的数据流
	d, err := io.ReadAll(response.Body)
	if err != nil {
		return p, err
	}
	return preload.Unmarshal(d, "json", encryptionKey)
}
