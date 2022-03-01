package giotgo

import (
	"fmt"
	"net"

	"github.com/ghuvrons/g-IoT-Go/giot_packet"
)

type Server struct {
	clients []*ClientHandler

	setting struct {
		clientTimeout int
	}

	authenticator    func(username, password string) bool
	commandHandlers  map[giot_packet.Command]CommandHandler
	commandExecutors map[giot_packet.Command]CommandExecutor
}

func NewServer() *Server {
	server := &Server{}
	server.commandHandlers = map[giot_packet.Command]CommandHandler{}
	server.commandExecutors = map[giot_packet.Command]CommandExecutor{}
	return server
}

func (svr *Server) Serve(addr string) {
	defer func() {
		fmt.Println("server closed")
	}()

	serverSock, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	for true {
		conn, err := serverSock.Accept()
		if err != nil {
			continue
		}
		fmt.Println("new client")
		client := NewClientHandler(&conn, svr, svr.setting.clientTimeout)
		svr.clients = append(svr.clients, client)
	}
}

func (svr *Server) On(cmd giot_packet.Command, handler CommandHandler) {
	svr.commandHandlers[cmd] = handler
}

func (svr *Server) OnExecute(cmd giot_packet.Command, handler CommandExecutor) {
	svr.commandExecutors[cmd] = handler
}

func (svr *Server) ClientAuth(authenticator func(username, password string) bool) {
	svr.authenticator = authenticator
}
