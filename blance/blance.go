package main

import (
	"info"
	"mclient"
	"mserver"
	"time"
)

type Blance struct {
	mserver *mserver.Manager
	mclient *mclient.Manager
	info    *info.Info
}

func NewBlance() *Blance {
	return &Blance{
		mserver: mserver.NewManager(),
		mclient: mclient.NewManager(),
		info:    info.NewInfo(),
	}
}

func (blance *Blance) Start() {
	for {
		select {
		case <-blance.info.HasClient:
			blance.mserver.Look()
			break
		}
	}
}

func (blance *Blance) Demon() {
	for {
		if blance.mserver.Listener == nil {
			blance.mserver.Start(blance.info)
		}
		if blance.mclient.Listener == nil {
			blance.mclient.Start(blance.info)
		}
		time.Sleep(time.Second)
	}
}

func (blance *Blance) Close() {
	blance.mserver.Close()
	blance.mclient.Close()
}
