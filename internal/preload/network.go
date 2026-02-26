package preload

import (
	"github.com/unix755/xtools/xNet"
)

func (p Preload) GetInterface(name string) xNet.NetInterface {
	for _, netInterface := range p.NetInterfaces {
		if netInterface.Name == name {
			return netInterface
		}
	}
	return xNet.NetInterface{}
}
