package main

import (
	"log"
	"os"
	"time"
)

const (
	BUF_LEN = 1024
)

var (
	host, port string      = "0.0.0.0", "8090"
	blance     *Blance     = nil
	loger      *log.Logger = nil
)

func init() {
	blance = NewBlance()
	loger = log.New(os.Stdout, "[BLANCE] ", log.Ldate|log.Ltime)
}

func Demon() {
	for {
		if blance.listener == nil {
			blance.Listen()
		}
		time.Sleep(time.Second)
	}
}

func main() {
	defer func() {
		blance.Close()
	}()

	//start blance
	blance.Start()
	Demon()
}
