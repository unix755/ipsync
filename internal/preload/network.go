package preload

import (
	"net/netip"

	"github.com/unix755/xtools/xNet"
)

// GetFirstPublicIP 获取负载中指定 remoteInterface 的第一个公网 IP
func (p Preload) GetFirstPublicIP(remoteInterface string) (netip.Addr, bool) {
	// 获取公共 IP
	for _, netInterface := range p.NetInterfaces {
		if netInterface.Name == remoteInterface {
			for _, ip := range netInterface.IPs {
				isPublic, err := xNet.IsPublic(ip.String())
				if err == nil && isPublic {
					return ip, true
				}
			}
		}
	}
	return netip.Addr{}, false
}
