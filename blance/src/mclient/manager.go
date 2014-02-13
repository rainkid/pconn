package mclient

import (
	"encoding/json"
	"fmt"
	"info"
	"log"
	"net"
	"os"
)

const (
	BUF_LEN   = 1024
	CMD_QUERY = 1001
)

var (
	host, port string      = "0.0.0.0", "7080"
	loger      *log.Logger = log.New(os.Stdout, "[blance-client-manager] ", log.Ldate|log.Ltime)
)

type Manager struct {
	Listener *net.TCPListener
	qclose   chan *Client
	qconnect chan *net.TCPConn
	qmessage chan *Message
	info     *info.Info
}

type Message struct {
	data   []byte
	client *Client
}

type Cmd struct {
	Code int
	Data []byte
}

func NewManager() *Manager {
	return &Manager{
		Listener: nil,
		qclose:   make(chan *Client),
		qconnect: make(chan *net.TCPConn),
		qmessage: make(chan *Message),
	}
}

func (cm *Manager) Start(info *info.Info) {
	cm.info = info

	err := cm.Listen()
	if err != nil {
		loger.Println("error connect : ", err.Error())
		return
	}
	//accept client
	go cm.AcceptTcp()
	//dispatch message
	go cm.Dispatch()
}

func (cm *Manager) AcceptTcp() {
	for {
		if cm.Listener != nil {
			conn, err := cm.Listener.AcceptTCP()
			if err != nil {
				loger.Println("error accept:", err.Error())
				continue
			}
			cm.info.HasClient <- true
			cm.qconnect <- conn
		}
	}
}

func (cm *Manager) Dispatch() {
	for {
		select {
		case client := <-cm.qclose:
			cm.DelClient(client)
			break
		case conn := <-cm.qconnect:
			cm.AddClient(conn)
			break
		case message := <-cm.qmessage:
			cm.Qmessage(message)
			break
		}
	}
}

func (cm *Manager) Qmessage(msg *Message) {
	client := msg.client

	var cmd Cmd

	err := json.Unmarshal(msg.data, &cmd)
	if err != nil {
		loger.Println("error unmarshal.", err.Error())
		return
	}
	fmt.Println(cmd.Code)
	fmt.Println(string(cmd.Data))
	//
	switch cmd.Code {
	case CMD_QUERY:
		str := fmt.Sprintf("%s:%d", cm.info.ServerHost, cm.info.ServerPort)
		loger.Println("dispatch info,", str)

		var m = Cmd{CMD_QUERY, []byte(str)}
		s, err := json.Marshal(m)
		if err != nil {
			loger.Println("send error,", err.Error())
			break
		}
		client.Send(s)
		break
	}
}

func (cm *Manager) Listen() error {
	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	loger.Println("started listen.")
	cm.Listener = listener
	return nil
}

func (cm *Manager) AddClient(conn *net.TCPConn) {
	client := NewClient(cm, conn)
	//client start
	client.Start()
	loger.Println(client.host, client.port, " connected.")
}

func (cm *Manager) DelClient(client *Client) {
	loger.Println(client.host, client.port, " closed.")
	client.Close()
}

func (cm *Manager) Close() {
	if cm.Listener != nil {
		cm.Listener.Close()
		cm.Listener = nil
	}
}
