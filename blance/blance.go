package main

import (
	"encoding/json"
	"net"
	"time"
)

const (
	CMD_PING = 1
)

var (
	frequency         int64 = 2
	heartTimeOutTimes int64 = 2
)

type Blance struct {
	servers  map[uint32]*Server
	listener *net.TCPListener
	qclose   chan *Server
	qconnect chan *net.TCPConn
	qmessage chan *Message
}

type Message struct {
	data   []byte
	server *Server
}

type Cmd struct {
	Code int
	Data []byte
}

func NewBlance() *Blance {
	return &Blance{
		servers:  make(map[uint32]*Server),
		listener: nil,
		qclose:   make(chan *Server),
		qconnect: make(chan *net.TCPConn),
		qmessage: make(chan *Message),
	}
}

//heart broad to all server
func (blance *Blance) HeartBroad() {
	for {
		for _, server := range blance.servers {

			cmd := &Cmd{CMD_PING, []byte("heart")}
			buf, err := json.Marshal(cmd)
			if err != nil {
				loger.Println("cmd json error.")
				continue
			}

			//send heart
			server.Send(buf)

			//heart time timeout
			ntime := time.Now().Unix()
			if (ntime - server.lastHeart) > (heartTimeOutTimes * frequency) {
				blance.DelServer(server)
			}
		}
		time.Sleep(time.Second * time.Duration(frequency))
	}
}

func (blance *Blance) Dispatch() {
	for {
		select {
		case server := <-blance.qclose:
			blance.DelServer(server)
			break
		case conn := <-blance.qconnect:
			blance.AddServer(conn)
			break
		case message := <-blance.qmessage:
			blance.Qmessage(message)
			break
		}
	}
}

func (blance *Blance) Start() {
	listener, err := blance.Listen()
	if err != nil {
		loger.Println("error connect : ", err.Error())
		return
	}
	blance.listener = listener

	//start accept
	go blance.AcceptTcp()
	//start heartbroad runtime
	go blance.HeartBroad()
	//dispatch conn msg
	go blance.Dispatch()
}

func (blance *Blance) Listen() (*net.TCPListener, error) {

	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		loger.Println("listen faild.")
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		loger.Println("listen faild.")
		return nil, err
	}
	loger.Println("started listen.")
	return listener, nil
}

func (blance *Blance) AcceptTcp() {
	for {
		if blance.listener != nil {
			conn, err := blance.listener.AcceptTCP()
			if err != nil {
				loger.Println("error accept:", err.Error())
				continue
			}

			blance.qconnect <- conn
		}
	}
}

func (blance *Blance) Qmessage(msg *Message) {
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
		loger.Println(server)
		server.stat = &stat
		break
	}
}

func (blance *Blance) AddServer(conn *net.TCPConn) {
	server := NewServer(conn)
	loger.Println(server.host, server.port, " connected.")
	if _, ok := blance.servers[server.uhost]; !ok {
		blance.servers[server.uhost] = server
	}
	//server loop and deal msg
	server.Start()
}

func (blance *Blance) DelServer(server *Server) {
	loger.Println(server.host, server.port, " closed.")
	delete(blance.servers, server.uhost)
	server.Close()
}

func (blance *Blance) Close() {
	if blance.listener != nil {
		blance.listener.Close()
		blance.listener = nil
	}

	for _, server := range blance.servers {
		server.Close()
	}
}
