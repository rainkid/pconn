package main

import (
	"runtime"
)

var (
	blance *Blance
)

func init() {
	blance = NewBlance()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	defer func() {
		blance.Close()
	}()

	go blance.Start()
	//start blance
	blance.Demon()
}
