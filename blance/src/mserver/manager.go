package mserver

import (
	"encoding/json"
	"errors"
	"info"
	"log"
	"net"
	"os"
	"time"
)

const (
	BUF_LEN  = 1024
	CMD_PING = 1001
)

var (
	host, port        string      = "0.0.0.0", "8090"
	frequency         int64       = 20
	heartTimeOutTimes int64       = 2
	loger             *log.Logger = log.New(os.Stdout, "[blance-server-manager] ", log.Ldate|log.Ltime)
)

type Manager struct {
	Servers  map[uint32]*Server
	Listener *net.TCPListener
	qclose   chan *Server
	qconnect chan *net.TCPConn
	qmessage chan *Message
	info     *info.Info
}

type Message struct {
	data   []byte
	server *Server
}

type Cmd struct {
	Code int
	Data []byte
}

func NewManager() *Manager {
	return &Manager{
		Servers:  make(map[uint32]*Server),
		Listener: nil,
		qclose:   make(chan *Server),
		qconnect: make(chan *net.TCPConn),
		qmessage: make(chan *Message),
	}
}

func (sm *Manager) Look() error {
	var weight = 100000
	var cserver *Server

	for _, server := range sm.Servers {
		sweight := int(server.weight)
		if sweight < weight {
			weight = sweight
			cserver = server
		}
	}
	if cserver == nil {
		return errors.New("no server")
	}

	sm.info.ServerHost, sm.info.ServerPort = cserver.host, 9010
	return nil
}

//heart broad to all server
func (sm *Manager) HeartBroad() {
	for {
		for _, server := range sm.Servers {
			server.Heart()
		}
		time.Sleep(time.Second * time.Duration(frequency))
	}
}

func (sm *Manager) Dispatch() {
	for {
		select {
		case server := <-sm.qclose:
			sm.DelServer(server)
			break
		case conn := <-sm.qconnect:
			sm.AddServer(conn)
			break
		case message := <-sm.qmessage:
			sm.Qmessage(message)
			break
		}
	}
}

func (sm *Manager) Start(info *info.Info) {
	sm.info = info

	err := sm.Listen()
	if err != nil {
		loger.Println("error connect : ", err.Error())
		return
	}

	//start accept
	go sm.AcceptTcp()
	//start heartbroad runtime
	go sm.HeartBroad()
	//dispatch conn msg
	go sm.Dispatch()
}

func (sm *Manager) Listen() error {
	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	loger.Println("started listen.")
	sm.Listener = listener
	return nil
}

func (sm *Manager) AcceptTcp() {
	for {
		if sm.Listener != nil {
			conn, err := sm.Listener.AcceptTCP()
			if err != nil {
				loger.Println("error accept:", err.Error())
				continue
			}

			sm.qconnect <- conn
		}
	}
}

func (sm *Manager) Qmessage(msg *Message) {
	server := msg.server
	server.lastHeart = time.Now().Unix()

	var cmd Cmd
	err := json.Unmarshal(msg.data, &cmd)
	if err != nil {
		return
	}
	//
	switch cmd.Code {
	case CMD_PING:
		var stat Stat
		err = json.Unmarshal(cmd.Data, &stat)
		if err != nil {
			loger.Println("ping unmarshal error:", err.Error())
			break
		}
		server.stat = &stat
		sm.Weight(server)
		loger.Println(server.host, server.weight)
		break
	}
}

func (sm *Manager) AddServer(conn *net.TCPConn) {
	server := NewServer(sm, conn)

	if _, ok := sm.Servers[server.uhost+server.port]; !ok {
		sm.Servers[server.uhost+server.port] = server
	}

	//server start
	server.Start()
	//server heart
	server.Heart()

	loger.Println(server.host, server.port, " connected, total:", len(sm.Servers))
}

func (sm *Manager) DelServer(server *Server) {
	loger.Println(server.host, server.port, " closed.")
	delete(sm.Servers, server.uhost)
	server.Close()
}

func (sm *Manager) Weight(server *Server) {
	stat := server.stat

	load15 := float64(stat.Load15) * 100
	memused := stat.Memused * 100

	server.weight = int64(int(load15) + int(memused) + stat.Clients)
}

func (sm *Manager) Close() {
	if sm.Listener != nil {
		sm.Listener.Close()
		sm.Listener = nil
	}

	for _, server := range sm.Servers {
		server.Close()
	}
}
