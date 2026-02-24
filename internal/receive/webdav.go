package receive

import (
	"io"
	"log"
	"netinfo/internal/preload"
	"netinfo/internal/wireguard"
	"time"

	"github.com/unix755/xtools/xWebDAV"
)

// getNetInfoFromWebDAV 从 webdav 服务器获取指定 id 的网络信息
func getNetInfoFromWebDAV(endpoint string, username string, password string, allowInsecure bool, filepath string, encryptionKey []byte) (p *preload.Preload, err error) {
	client, err := xWebDAV.NewClient(endpoint, username, password, allowInsecure)
	if err != nil {
		return nil, err
	}
	response, err := client.Download(filepath)
	if err != nil {
		return nil, err
	}

	// 读取从 webdav 服务器下载的数据流
	d, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return preload.Unmarshal(d, "json", encryptionKey)
}

func FromWebDAV(endpoint string, username string, password string, allowInsecure bool, filepath string, encryptionKey []byte, remoteInterface string, wgInterface string, wgPeerKey string) (err error) {
	p, err := getNetInfoFromWebDAV(endpoint, username, password, allowInsecure, filepath, encryptionKey)
	if err != nil {
		return err
	}
	publicIP, err := p.GetPublicIP(remoteInterface)
	if err != nil {
		return err
	}
	return wireguard.UpdateEndpoint(wgInterface, wgPeerKey, publicIP, -1)
}

func FromWebDAVLoop(endpoint string, username string, password string, allowInsecure bool, filepath string, encryptionKey []byte, remoteInterface string, wgInterface string, wgPeerKey string, interval time.Duration) {
	for {
		err := FromWebDAV(endpoint, username, password, allowInsecure, filepath, encryptionKey, remoteInterface, wgInterface, wgPeerKey)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(interval)
	}
}
