// testdata/10.获取地理位置.go
package main

import (
	"blogx_server/core"
	"fmt"
)

func main() {
	ip2region()
}
func ip2region() {
	core.InitIPDB()

	ips := []string{
		"172.16.22.1",     // 内网 IP
		"127.0.0.1",       // 本地回环
		"175.0.201.207",   // 广东深圳
		"123.123.123.123", // 江苏南京
		"8.8.8.8",         // 美国 Google DNS
		"203.198.23.69",   // 香港
		"61.216.152.102",  // 台湾
		"abc.def.ghi",     // 非法 IP
		"",                // 空字符串
		"999.999.999.999", // 边界非法 IP
	}

	for _, ip := range ips {
		fmt.Printf("IP: %-16s => %s\n", ip, core.GetIpAddr(ip))
	}
}
