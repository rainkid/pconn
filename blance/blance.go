package main

import (
	"fmt"
	"net"
)

type Blance struct {
	pings    chan *Center
	centers  map[uint32]*Center
	listener *net.TCPListener
}

func (blance *Blance) Dispatch() {
	for {
		select {
		case center := <-blance.pings:
			//send reply
			_, err := center.Write(CMD_OK)
			if err != nil {
				fmt.Println("error send reply:", err.Error())
				center.Close()
			}
			break
		}
	}
}

func (blance *Blance) Start() {
	listener, err := blance.Listen()
	if err != nil {
		fmt.Println("error connect : ", err.Error())
		return
	}
	blance.listener = listener

	//start accept
	go blance.AcceptTcp()
	//dispatch conn msg
	go blance.Dispatch()
}

func (blance *Blance) Listen() (*net.TCPListener, error) {
	fmt.Print("start listened")

	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		fmt.Println("...faild")
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("...faild")
		return nil, err
	}
	fmt.Println("...success.")
	return listener, nil
}

func (blance *Blance) AcceptTcp() {
	for {
		if blance.listener != nil {
			conn, err := blance.listener.AcceptTCP()
			if err != nil {
				fmt.Println("error accept:", err.Error())
				continue
			}
			blance.NewCenter(conn)
		}
	}
}

func (blance *Blance) NewCenter(conn *net.TCPConn) {
	center := NewCenter(conn)
	if _, ok := blance.centers[center.uhost]; !ok {
		blance.centers[center.uhost] = center
	}
	//center loop and deal msg
	go center.Demon()
}

func (blance *Blance) Close() {
	if blance.listener != nil {
		blance.listener.Close()
		blance.listener = nil
	}

	for _, center := range blance.centers {
		center.Close()
	}
}
