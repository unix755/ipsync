package preload

import (
	"ipsync/internal/cache"
	"ipsync/internal/wireguard"
	"log"
	"time"
)

func (p *Preload) UpdateWireGuardEndPoint(remoteInterface string, wgInterface string, wgPeerKey string, newEndpointPort int, interval time.Duration) (err error) {
	if interval != 0 {
		for {
			// 获取第一个公网 IP
			ip, found := p.GetFirstPublicIP(remoteInterface)
			if found {
				// 获取缓存中的 wg_endpoint
				cacheWGEndpoint, _ := cache.Get("wg_endpoint")

				if ip.String() != cacheWGEndpoint {
					// 更新 wireguard 端点
					err = wireguard.UpdateEndpoint(wgInterface, wgPeerKey, ip.String(), newEndpointPort)
					if err != nil {
						log.Println(err)
					} else {
						// 设置新的缓存 wg_endpoint
						cache.Set("wg_endpoint", ip.String())
						log.Println("new wireguard endpoint ip found, updated wireguard endpoint to ", ip.String())
					}
				} else {
					log.Println("new wireguard endpoint ip not found, skip")
				}
			} else {
				log.Println("cannot find wireguard endpoint from remote interface")
			}

			time.Sleep(interval)
		}
	} else {
		// 获取到第一个公网 IP
		ip, found := p.GetFirstPublicIP(remoteInterface)

		if found {
			err = wireguard.UpdateEndpoint(wgInterface, wgPeerKey, ip.String(), newEndpointPort)
		}
	}
	return err
}
