package mserver

import (
	"encoding/json"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
	"utils"
)

type Server struct {
	manager   *Manager
	host      string
	uhost     uint32
	port      uint32
	conn      *net.TCPConn
	lastHeart int64
	stat      *Stat
	weight    int64
}

type Stat struct {
	Load1   float64
	Load5   float64
	Load15  float64
	Cpunum  int
	Memused float64
	Clients int
}

func NewServer(manager *Manager, conn *net.TCPConn) *Server {
	ip := strings.Split(conn.RemoteAddr().String(), ":")
	port, _ := strconv.ParseUint(ip[1], 10, 64)

	return &Server{
		manager:   manager,
		host:      ip[0],
		uhost:     utils.Ip2Uint32(ip[0]),
		port:      uint32(port),
		conn:      conn,
		lastHeart: time.Now().Unix(),
		stat:      nil,
	}
}

func (server *Server) Heart() {
	cmd := &Cmd{CMD_PING, []byte("heart")}
	buf, err := json.Marshal(cmd)
	if err != nil {
		loger.Println("cmd json error.")
		return
	}

	//send heart
	_, err = server.Send(buf)
	if err != nil {
		loger.Println("send error,", err.Error())
	}

	//heart time timeout
	ntime := time.Now().Unix()
	if (ntime - server.lastHeart) > (heartTimeOutTimes * frequency) {
		server.manager.qclose <- server
	}
	return
}

func (server *Server) Write(s string) (int, error) {
	n, err := server.conn.Write([]byte(s))
	return n, err
}

func (server *Server) Read(buf []byte) (int, error) {
	n, err := server.conn.Read(buf)
	return n, err
}

func (server *Server) Send(buf []byte) (int, error) {
	if server.conn != nil {
		n, err := server.Write(string(buf))
		if err != nil {
			server.manager.qclose <- server
			loger.Println(server.host, server.port, "send error : ", err.Error())
			return n, err
		}
	}
	return 0, nil
}

func (server *Server) Listen() {
	for {
		buf := make([]byte, BUF_LEN)
		n, err := server.Read(buf)

		if err == io.EOF {
			server.manager.qclose <- server
			break
		}

		if err != nil {
			loger.Println("read error : ", err.Error())
			continue
		}

		msg := &Message{data: buf[0:n], server: server}
		server.manager.qmessage <- msg
	}
}

func (server *Server) Start() {
	go server.Listen()
}

func (server *Server) Close() {
	server.conn.Close()
	server.conn = nil
}
