package preload

import (
	"fmt"

	"github.com/unix755/xtools/xNet"
)

// GetPublicIP 从网络信息中获取公共 IP
func (p *Preload) GetPublicIP(interfaceName string) (ip string, err error) {
	for _, netInterface := range p.NetInterfaces {
		if netInterface.Name == interfaceName {
			for _, ip := range netInterface.IPs {
				isPublic, _ := xNet.IsPublic(ip.String())
				if isPublic {
					return ip.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("no valid public IP found in network infomation data")
}
