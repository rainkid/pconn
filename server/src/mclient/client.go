package mclient

import (
	"net"
	"strings"
	"utils"
)

type Client struct {
	manager *Manager
	host    string
	uhost   uint32
	port    string
	conn    *net.TCPConn
}

func NewClient(manager *Manager, conn *net.TCPConn) *Client {
	ip := strings.Split(conn.RemoteAddr().String(), ":")
	return &Client{
		manager: manager,
		host:    ip[0],
		uhost:   utils.Ip2Uint32(ip[0]),
		port:    ip[1],
		conn:    conn,
	}
}

func (client *Client) Write(s string) (int, error) {
	n, err := client.conn.Write([]byte(s))
	return n, err
}

func (client *Client) Read(buf []byte) (int, error) {
	n, err := client.conn.Read(buf)
	return n, err
}

func (client *Client) Send(buf []byte) (int, error) {
	if client.conn != nil {
		n, err := client.Write(string(buf))
		if err != nil {
			loger.Println(client.host, client.port, "send error : ", err.Error())
			return n, err
		}
	}
	return 0, nil
}

func (client *Client) Close() {
	client.conn.Close()
	client.conn = nil
}
