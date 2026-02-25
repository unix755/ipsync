package send

import (
	"encoding/json"
	"ipsync/internal/cache"
	"ipsync/internal/preload"
	"log"
	"net/http"
	"time"

	"github.com/unix755/xtools/xWebDAV"
)

func ToWebDAV(endpoint string, username string, password string, allowInsecure bool, filepath string, encryptionKey []byte) (resp *http.Response, err error) {
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

	client, err := xWebDAV.NewClient(endpoint, username, password, allowInsecure)
	if err != nil {
		return nil, err
	}

	return client.Upload(filepath, bytes)
}

func ToWebDAVLoop(endpoint string, username string, password string, allowInsecure bool, filepath string, encryptionKey []byte, interval time.Duration) {
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
			resp, err := ToWebDAV(endpoint, username, password, allowInsecure, filepath, encryptionKey)
			if err != nil {
				log.Println(err)
			} else {
				// 设置新的缓存 net_interfaces
				cache.Set("net_interfaces", string(bytes))
				log.Printf("new net interfaces found, response %s", resp.Status)
			}
		} else {
			log.Println("new net interfaces not found, skip")
		}

		time.Sleep(interval)
	}
}
