package send

import (
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
		resp, err := ToWebDAV(endpoint, username, password, allowInsecure, filepath, encryptionKey)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(resp.Status)
		}
		time.Sleep(interval)
	}
}
