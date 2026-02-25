package send

import (
	"encoding/json"
	"ipsync/internal/cache"
	"ipsync/internal/preload"
	"log"
	"os"
	"time"
)

func ToFile(file string, encryptionKey []byte) (err error) {
	// 获取负载
	p, err := preload.NewPreload()
	if err != nil {
		return err
	}
	// 负载转换为加密比特流
	bytes, err := preload.Marshal(p, "json", encryptionKey)
	if err != nil {
		return err
	}
	return os.WriteFile(file, bytes, 0644)
}

func ToFileLoop(file string, encryptionKey []byte, interval time.Duration) {
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
			err = ToFile(file, encryptionKey)
			if err != nil {
				log.Println(err)
			} else {
				// 设置新的缓存 net_interfaces
				cache.Set("net_interfaces", string(bytes))
				log.Printf("new net interfaces found, save file to %s\n", file)
			}
		} else {
			log.Println("new net interfaces not found, skip")
		}

		time.Sleep(interval)
	}
}
