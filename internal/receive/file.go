package receive

import (
	"log"
	"netinfo/internal/preload"
	"netinfo/internal/wireguard"
	"os"
	"time"
)

// getNetInfoFromFile 从 file 文件获取指定 id 的网络信息
func getNetInfoFromFile(filepath string, encryptionKey []byte) (p *preload.Preload, err error) {
	d, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return preload.Unmarshal(d, "json", encryptionKey)
}

func FromFile(filepath string, encryptionKey []byte, remoteInterface string, wgInterface string, wgPeerKey string) (err error) {
	p, err := getNetInfoFromFile(filepath, encryptionKey)
	if err != nil {
		return err
	}
	publicIP, err := p.GetPublicIP(remoteInterface)
	if err != nil {
		return err
	}
	return wireguard.UpdateEndpoint(wgInterface, wgPeerKey, publicIP, -1)
}

func FromFileLoop(filepath string, encryptionKey []byte, remoteInterface string, wgInterface string, wgPeerKey string, interval time.Duration) {
	for {
		err := FromFile(filepath, encryptionKey, remoteInterface, wgInterface, wgPeerKey)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(interval)
	}
}
