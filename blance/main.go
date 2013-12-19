package main

import (
	"time"
)

const (
	BUF_LEN      = 1024
	CMD_PING     = "cmd_ping"
	CMD_OK       = "cmd_ok"
	CMD_REGISTER = "cmd_register"
)

var (
	host   = "0.0.0.0"
	port   = "6666"
	blance *Blance
)

func init() {
	blance = &Blance{make(chan *Center), make(map[uint32]*Center), nil}
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
