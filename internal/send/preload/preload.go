package preload

import (
	"encoding/json"
	"netinfo/internal/netinfo"
	"time"

	"github.com/unix755/xtools/xCrypto"
	"golang.org/x/crypto/chacha20poly1305"
)

func newPreload() (preload []byte, err error) {
	var preloadStrut netinfo.Data

	netInterfaces, err := netinfo.GetNetInterfaces()
	if err != nil {
		return nil, err
	}

	preloadStrut.UpdatedAt = time.Now()
	preloadStrut.NetInterfaces = netInterfaces

	return json.Marshal(preloadStrut)
}

func newEncryptedPreload(key []byte) (preload []byte, err error) {
	plaintext, err := newPreload()
	if err != nil {
		return nil, err
	}
	return xCrypto.NewChaCha20Poly1305(key, []byte{}).Encrypt(plaintext)
}

func GetPreload(key []byte) (preload []byte, err error) {
	// 通过密钥长度判断是否使用加密
	switch len(key) {
	case 0:
		return newPreload()
	default:
		key = xCrypto.ZeroPadding(key, chacha20poly1305.KeySize)
		key = key[0:chacha20poly1305.KeySize]
		return newEncryptedPreload(key)
	}
}
