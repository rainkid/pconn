package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

var (
	client *Client
)

type Cmd struct {
	Code int
	Data []byte
}

type Client struct {
	conn       *net.TCPConn
	ServerInfo string
	message    chan []byte
}

func init() {
	client = &Client{
		conn:    nil,
		message: make(chan []byte),
	}
}

func (client *Client) Connect(host, port string) (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (client *Client) Write(s string) (int, error) {
	n, err := client.conn.Write([]byte(s))
	return n, err
}

func (client *Client) Read(buf []byte) (int, error) {
	n, err := client.conn.Read(buf)
	return n, err
}

func (client *Client) Listen() {
	for {
		if client.conn != nil {
			data := make([]byte, 1024)
			n, err := client.conn.Read(data)
			if err == io.EOF {
				client.conn.Close()
				continue
			}

			//send message
			client.message <- data[0:n]
		} else {
			break
		}
	}
}

func (client *Client) Message() {
	for {
		select {
		case msg := <-client.message:
			var cmd Cmd
			err := json.Unmarshal(msg, &cmd)
			if err != nil {
				fmt.Println("error unmarshal.", err.Error())
				break
			}

			client.ServerInfo = string(cmd.Data)
			info := strings.Split(string(cmd.Data), ":")
			conn, err := client.Connect(info[0], info[1])
			if err != nil {
				fmt.Println("message connect error,", err.Error())
				break
			}

			fmt.Println("start listen")
			client.conn = conn
			go client.Listen()
			break
		}
	}
}

func (client *Client) Start() {
	conn, err := client.Connect("127.0.0.1", "7080")
	if err != nil {
		fmt.Println("connect error,", err.Error())
		return
	}

	cmd := &Cmd{1001, []byte("client test connect.")}
	buf, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println("cmd json error.")
		return
	}

	client.conn = conn
	go client.Listen()

	client.conn.Write(buf)

	go client.Message()
}

func (client *Client) Demon() {
	for {
		if client.ServerInfo == "" {
			client.Start()
		}
		time.Sleep(time.Second * 5)
	}
}

func main() {
	client.Demon()
}
