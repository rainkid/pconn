package main

import (
	"io"
	"net"
	"strings"
	"time"
	"utils"
)

type Server struct {
	host      string
	uhost     uint32
	port      string
	conn      *net.TCPConn
	lastHeart int64
	stat      *Stat
}

type Stat struct {
	Load1   float64
	Load5   float64
	Load15  float64
	Cpunum  int
	Memused float64
}

func NewServer(conn *net.TCPConn) *Server {
	ip := strings.Split(conn.RemoteAddr().String(), ":")
	return &Server{
		host:      ip[0],
		uhost:     utils.Ip2Uint32(ip[0]),
		port:      ip[1],
		conn:      conn,
		lastHeart: time.Now().Unix(),
		stat:      nil,
	}
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
			blance.qclose <- server
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
			blance.qclose <- server
			break
		}

		if err != nil {
			loger.Println("read error : ", err.Error())
			continue
		}

		msg := &Message{data: buf[0:n], server: server}
		blance.qmessage <- msg
	}
}

func (server *Server) Start() {
	go server.Listen()
}

func (server *Server) Close() {
	server.conn.Close()
	server.conn = nil
}
