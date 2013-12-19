package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

const (
	BUF_LEN  = 1024
	CMD_PING = "cmd_ping"
	CMD_OK   = "cmd_ok"
)

var (
	host   = "0.0.0.0"
	port   = "6666"
	global *Global
)

type Global struct {
	relisten chan bool
	conn     *net.TCPConn
}

func init() {
	global = &Global{make(chan bool), nil}
}

func Demon() {
	for {
		if global.conn == nil {
			Listen()
		}
		time.Sleep(time.Second)
	}
}

func Ping() {
	for {
		if global.conn != nil {
			_, err := Write(global.conn, CMD_PING)
			if err != nil {
				fmt.Println("heartbeat error : ", err.Error())
			}

			data, err := Read(global.conn)
			if err == io.EOF {
				fmt.Println("heartbeat read error:", err.Error())
				Destory()
			}
			fmt.Println("heartbeat recived:", string(data))
		}
		time.Sleep(time.Second)
	}
}

func Write(conn *net.TCPConn, s string) (int, error) {
	n, err := conn.Write([]byte(s))
	return n, err
}

func Read(conn *net.TCPConn) ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	fmt.Println("received ", n, " bytes of data =", string(buf[0:n]))
	return buf, err
}

func Listen() {
	conn, err := Connect()
	if err != nil {
		fmt.Println("error connect : ", err.Error())
		return
	}
	global.conn = conn
}

func Connect() (*net.TCPConn, error) {
	fmt.Println("try connect server.")
	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Destory() {
	if global.conn != nil {
		global.conn.Close()
		global.conn = nil
	}
}

func main() {
	defer func() {
		Destory()
	}()
	go Ping()
	Demon()
}
