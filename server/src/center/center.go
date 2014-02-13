package center

import (
	"flag"
	"mclient"
	"server"
	"time"
)

var (
	host string
	port string
)

type Center struct {
	server   *server.Server
	cmanager *mclient.Manager
}

func NewCenter() *Center {
	flag.StringVar(&host, "host", "127.0.0.1", "blance host.")
	flag.StringVar(&port, "port", "8090", "blance port.")

	return &Center{
		server:   server.NewServer(),
		cmanager: mclient.NewManager(),
	}
}

func (center *Center) Demon() {
	flag.Parse()

	for {
		if center.server.Conn == nil {
			center.server.Start(host, port)
		}
		if center.cmanager.Listener == nil {
			center.cmanager.Start()
		}
		time.Sleep(time.Second * time.Duration(5))
	}
}

func (center *Center) Close() {
	center.server.Close()
	center.cmanager.Close()
}
