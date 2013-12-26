package main

import (
	"net"
)

type Client struct {
	conn *net.TCPConn
}
