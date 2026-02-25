package preload

import (
	"encoding/json"
	"encoding/xml"
	"ipsync/internal/network"
	"time"
)

type Preload struct {
	UpdatedAt     time.Time              `json:"updatedAt,omitempty" xml:"updatedAt,omitempty" form:"updatedAt,omitempty"`
	NetInterfaces []network.NetInterface `json:"netInterfaces" xml:"netInterfaces" form:"netInterfaces" binding:"required"`
}

func NewPreload() (preload Preload, err error) {
	netInterfaces, err := network.GetNetInterfaces()
	if err != nil {
		return preload, err
	}

	return Preload{
		UpdatedAt:     time.Now(),
		NetInterfaces: netInterfaces,
	}, nil
}

func Marshal(preload Preload, preloadType string, key []byte) (preloadBytes []byte, err error) {
	// preload 转换为比特流
	switch preloadType {
	case "json":
		preloadBytes, err = json.Marshal(preload)
	case "xml":
		preloadBytes, err = xml.Marshal(preload)
	}
	if err != nil {
		return nil, err
	}

	// 比特流加密
	return Encrypt(preloadBytes, key)
}

func Unmarshal(preloadBytes []byte, preloadType string, key []byte) (preload *Preload, err error) {
	// 比特流解密
	preloadBytes, err = Decrypt(preloadBytes, key)
	if err != nil {
		return nil, err
	}

	// 比特流转换为 preload
	switch preloadType {
	case "json":
		err = json.Unmarshal(preloadBytes, preload)
	case "xml":
		err = xml.Unmarshal(preloadBytes, preload)
	}
	return preload, err
}
