package main

import (
	"center"
	"runtime"
)

var (
	mcenter *center.Center
)

func init() {
	mcenter = center.NewCenter()
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	defer func() {
		mcenter.Close()
	}()

	mcenter.Demon()
}
