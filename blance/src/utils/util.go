package utils

import (
	"fmt"
)

func Ip2Uint32(str string) uint32 {
	var a, b, c, d byte
	n, err := fmt.Sscanf(str, "%d.%d.%d.%d", &a, &b, &c, &d)
	if err != nil || n != 4 {
		return 0
	}
	ip := uint32(a)
	ip |= uint32(b) << 8
	ip |= uint32(c) << 16
	ip |= uint32(d) << 24
	return ip
}

func Ip2String(ip uint32) string {
	a := byte(ip & 0xff)
	b := byte((ip >> 8) & 0xff)
	c := byte((ip >> 16) & 0xff)
	d := byte((ip >> 24) & 0xff)
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}
