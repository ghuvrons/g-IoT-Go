package giotgo

import (
	"fmt"
	"net"
)

type Server struct {
	clients []*ClientHandler

	setting struct {
		clientTimeout int
	}
}

func (svr *Server) Serve(addr string) {
	defer func() {
		fmt.Println("server closed")
	}()

	serverSock, _ := net.Listen("tcp", addr)

	for true {
		conn, err := serverSock.Accept()
		if err != nil {
			continue
		}
		fmt.Println("new client")
		client := NewClientHandler(&conn, svr.setting.clientTimeout)
		svr.clients = append(svr.clients, client)
	}
}
