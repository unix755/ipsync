package send

import (
	"log"
	"netinfo/internal/preload"
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
		err := ToFile(file, encryptionKey)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("save file to %s\n", file)
		}
		time.Sleep(interval)
	}
}
