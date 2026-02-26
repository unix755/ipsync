package wireguard

import (
	"fmt"
	"ipsync/internal/cache"
	"log"
	"strconv"
	"time"
)

// UpdateEndpoint 更新 endpoint
func UpdateEndpoint(wgInterface string, wgPeerKey string, newEndpointAddr string, newEndpointPort int) (err error) {
	// 指定了 peer key 就只更新匹配的 EndpointConfig, 未指定则更新所有的 EndpointConfig
	if wgPeerKey != "" {
		// 获取单个 EndpointConfig
		endpointConfig, err := GetEndpointConfigByKey(wgInterface, wgPeerKey)
		if err != nil {
			return fmt.Errorf("unable to get wireguard endpoints from interface: %s key: %s, error: %v, maybe Permission denied", wgInterface, wgPeerKey, err)
		}
		err = endpointConfig.ApplyNewEndpointConfig(wgInterface, wgPeerKey, newEndpointAddr, newEndpointPort)
		if err != nil {
			return err
		}
	} else {
		// 获取所有 EndpointConfig
		endpointConfigs, err := GetEndpointConfigs(wgInterface)
		if err != nil {
			return fmt.Errorf("unable to get wireguard endpoints from interface: %s, error: %v, maybe Permission denied", wgInterface, err)
		}
		for _, endpointConfig := range endpointConfigs {
			err = endpointConfig.ApplyNewEndpointConfig(wgInterface, wgPeerKey, newEndpointAddr, newEndpointPort)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateEndPointLoop 循环更新 endpoint
func UpdateEndPointLoop(wgInterface string, wgPeerKey string, newEndpointAddr string, newEndpointPort int, interval time.Duration) (err error) {
	if interval != 0 {
		for {
			// 缓存中查询
			wgInfo := wgInterface + wgPeerKey + newEndpointAddr + strconv.Itoa(newEndpointPort)
			cacheWGInfo, _ := cache.Get("wg_info")

			// 查询到变化即更新
			if wgInfo != cacheWGInfo {
				err = UpdateEndpoint(wgInterface, wgPeerKey, newEndpointAddr, newEndpointPort)
				if err != nil {
					log.Println(err)
				}
				cache.Set("wg_info", wgInfo)
			}

			// 等待下一次运行
			time.Sleep(interval)
		}
	} else {
		return UpdateEndpoint(wgInterface, wgPeerKey, newEndpointAddr, newEndpointPort)
	}
}
