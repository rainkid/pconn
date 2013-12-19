package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
	"utils"
)

type Center struct {
	host   string
	uhost  uint32
	port   string
	status int
	conn   *net.TCPConn
}

func NewCenter(conn *net.TCPConn) *Center {
	ip := strings.Split(conn.RemoteAddr().String(), ":")
	conn.SetReadDeadline(time.Time{})
	return &Center{ip[0], utils.Ip2Uint32(host), ip[1], 1, conn}
}

func (center *Center) Write(s string) (int, error) {
	n, err := center.conn.Write([]byte(s))
	return n, err
}

func (center *Center) Read(buf []byte) ([]byte, error) {
	_, err := center.conn.Read(buf)
	return buf, err
}

func (center *Center) Demon() {
	for {
		buf := make([]byte, 1024)
		data, err := center.Read(buf)

		if err == io.EOF {
			center.Close()
			break
		}

		if err != nil {
			fmt.Println("read error : ", err.Error())
			continue
		}
		fmt.Println("reviced from", center.host, center.port, "data =", string(data))
		center.Message(data)
	}
}

//center message center
func (center *Center) Message(data []byte) {
	blance.pings <- center
}

func (center *Center) Close() {
	center.conn.Close()
	fmt.Println(center.host, ":", center.port, " closed .")
}
