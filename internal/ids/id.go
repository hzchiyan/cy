package ids

import (
	"net"
	"strings"
)

func MacID() string {
	addrs, _ := net.Interfaces()
	var macs []string
	for _, address := range addrs {
		mac := address.HardwareAddr //获取本机MAC地址
		macs = append(macs, mac.String())
	}
	return strings.ReplaceAll(strings.Join(macs, ""), ":", "")
}
