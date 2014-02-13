package utils

import (
	"fmt"
	"runtime"
	"stat"
	"strconv"
)

type Stat struct {
	Load1   float64
	Load5   float64
	Load15  float64
	Cpunum  int
	Memused float64
	Clients int
}

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

func SysStat() *Stat {
	load := stat.GetLoadAvgSample()
	mem := stat.GetMemSample()

	memused := float64(mem.MemTotal - mem.MemFree - mem.Cached - mem.Buffers)
	memtotal := float64(mem.MemTotal)
	usedPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", memused/memtotal), 64)

	return &Stat{
		Load1:   load.One,
		Load5:   load.Five,
		Load15:  load.Fifteen,
		Cpunum:  runtime.NumCPU(),
		Memused: usedPercent,
	}
}
