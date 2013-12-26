package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"stat"
	"strconv"
	// "strings"
	"time"
)

const BUF_LEN = 1024
const (
	CMD_PING = 1
)

var (
	host      *string     = flag.String("host", "0.0.0.0", "tcp listen host")
	port      *string     = flag.String("port", "8090", "tcp listen port")
	server    *Server     = nil
	frequency int         = 60
	loger     *log.Logger = nil
)

type Server struct {
	conn    *net.TCPConn
	stat    *Stat
	Clients *Client
}

type Stat struct {
	Load1   float64
	Load5   float64
	Load15  float64
	Cpunum  int
	Memused float64
}

type Cmd struct {
	Code int
	Data []byte
}

func init() {
	server = &Server{
		conn: nil,
		stat: nil,
	}
	loger = log.New(os.Stdout, "[CENTER] ", log.Ldate|log.Ltime)
}

func Demon() {
	for {
		if server.conn == nil {
			server.Start()
		}
		time.Sleep(time.Second)
	}
}

func (server *Server) Ping() {
	for {
		if server.conn != nil {
			_, err := server.Write("stat")
			if err != nil {
				loger.Println("heartbeat error : ", err.Error())
			}

			data := make([]byte, BUF_LEN)
			_, err = server.Read(data)
			if err == io.EOF {
				loger.Println("heartbeat read error:", err.Error())
				server.Close()
			}
			loger.Println("heartbeat recived:", string(data))
		}
		time.Sleep(time.Second * 1)
	}
}

func (server *Server) Send(buf []byte) (int, error) {
	if server.conn != nil {
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
		if server.conn != nil {

			data := make([]byte, BUF_LEN)
			n, err := server.Read(data)
			if err == io.EOF {
				loger.Println("read error:", err.Error())
				server.Close()
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
		fmt.Println("cmd json unmarshal error:", err.Error())
	}

	//dispatch cmd
	switch cmd.Code {
	case CMD_PING:
		stat, _ := json.Marshal(server.Stat())
		ret := &Cmd{CMD_PING, stat}
		buf, _ := json.Marshal(ret)

		fmt.Println(string(buf))
		server.Send(buf)
		break
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

func (server *Server) Start() {
	conn, err := server.Connect()
	if err != nil {
		loger.Println("error connect : ", err.Error())
		return
	}
	server.conn = conn
}

func (server *Server) Connect() (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp", *host+":"+*port)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	loger.Println("connected server ", *host, ":", *port)
	return conn, nil
}

func (server *Server) Stat() *Stat {
	load := stat.GetLoadAvgSample()
	mem := stat.GetMemSample()

	memused := float64(mem.MemTotal - mem.MemFree - mem.Cached - mem.Buffers)
	memtotal := float64(mem.MemTotal)
	usedPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", memused/memtotal), 64)

	return &Stat{
		Load1:   load.One,
		Load5:   load.Five,
		Load15:  load.Fifteen,
		Cpunum:  runtime.NumCPU(),
		Memused: usedPercent,
	}
}

func (server *Server) Close() {
	if server.conn != nil {
		server.conn.Close()
		server.conn = nil
	}
}

func main() {
	flag.Parse()

	defer func() {
		server.Close()
	}()
	server.Stat()
	go server.Listen()
	Demon()
}
