package server

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"time"
	"utils"
)

const (
	BUF_LEN  = 1024
	CMD_PING = 1001
)

var (
	loger *log.Logger = log.New(os.Stdout, "[server] ", log.Ldate|log.Ltime)
)

type Server struct {
	Conn *net.TCPConn
	stat *utils.Stat
}

type Cmd struct {
	Code int
	Data []byte
}

func NewServer() *Server {
	return &Server{
		stat: utils.SysStat(),
	}
}

func (server *Server) Send(buf []byte) (int, error) {
	if server.Conn != nil {
		n, err := server.Write(string(buf))
		if err != nil {
			loger.Println("send error : ", err.Error())
			return n, err
		}
	}
	return 0, nil
}

func (server *Server) Listen() {
	for {
		if server.Conn != nil {

			data := make([]byte, BUF_LEN)
			n, err := server.Read(data)
			if err == io.EOF {
				server.Close()
				continue
			}

			//send message
			go server.Cmd(data[0:n])
		}
		time.Sleep(time.Second * 1)
	}
}

func (server *Server) Cmd(data []byte) {
	var cmd Cmd
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		loger.Println("cmd json unmarshal error:", err.Error())
	}

	//dispatch cmd
	switch cmd.Code {
	case CMD_PING:
		stat, _ := json.Marshal(utils.SysStat())
		ret := &Cmd{CMD_PING, stat}
		buf, _ := json.Marshal(ret)

		_, err := server.Send(buf)
		if err != nil {
			loger.Println("send error,", err.Error())
		}
		break
	}
}

func (server *Server) Write(s string) (int, error) {
	n, err := server.Conn.Write([]byte(s))
	return n, err
}

func (server *Server) Read(buf []byte) (int, error) {
	n, err := server.Conn.Read(buf)
	return n, err
}

func (server *Server) Start(host string, port string) {
	Conn, err := server.Connect(host, port)
	if err != nil {
		loger.Println("error Connect : ", err.Error())
		return
	}
	server.Conn = Conn

	go server.Listen()
}

func (server *Server) Connect(host string, port string) (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}

	Conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	loger.Println("Connected server ", host, ":", port)
	return Conn, nil
}

func (server *Server) Close() {
	loger.Println("server closed.")
	if server.Conn != nil {
		server.Conn.Close()
		server.Conn = nil
	}
}
