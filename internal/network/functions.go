package network

import (
	"net/netip"

	"github.com/unix755/xtools/xNet"
)

func GetNetInterfaces() (netInterfaces []NetInterface, err error) {
	nis, err := xNet.GetNetInterfaces()
	if err != nil {
		return nil, err
	}

	for _, ni := range nis {
		// 只取网路接口标记为UP的
		if !ni.Flag.Up {
			continue
		}

		// 拼接转换网络接口中的IPV4 与IPV6地址
		var ips []netip.Addr
		for _, ipString := range append(ni.Ipv4, ni.Ipv6...) {
			// 回环地址
			isLoopback, _ := xNet.IsLoopback(ipString)
			if isLoopback {
				continue
			}
			// 链路本地地址
			isLinkLocal, _ := xNet.IsLinkLocal(ipString)
			if isLinkLocal {
				continue
			}
			// 专用网络地址
			isPrivate, _ := xNet.IsPrivate(ipString)
			if isPrivate {
				//continue
			}
			// 地址转换出错
			ipAddr, err := netip.ParseAddr(ipString)
			if err != nil {
				continue
			}
			ips = append(ips, ipAddr)
		}

		// 跳过回环网络接口
		if len(ips) > 0 {
			netInterfaces = append(netInterfaces, NetInterface{
				Name: ni.Name,
				IPs:  ips,
				Mac:  ni.Mac,
			})
		}
	}

	return netInterfaces, nil
}
