package watchdog

import (
	"ipsync/internal/preload"
	"os"
)

// FromFile 从 file 文件获取指定 id 的网络信息
func FromFile(filepath string, encryptionKey []byte) (p preload.Preload, err error) {
	d, err := os.ReadFile(filepath)
	if err != nil {
		return p, err
	}
	return preload.Unmarshal(d, "json", encryptionKey)
}
