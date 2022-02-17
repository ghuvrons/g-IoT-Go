package main

import (
	"bytes"
	"fmt"

	giotgo "github.com/ghuvrons/g-IoT-Go"
	giot_packet "github.com/ghuvrons/g-IoT-Go/giot_packet"
)

func main() {
	fmt.Println("Start")

	var server = giotgo.NewServer()

	server.On(0xAAFF, func(client *giotgo.ClientHandler, data giot_packet.Data) *bytes.Buffer {
		fmt.Println(data)
		client.Execute(0xFFAA, data)
		return nil
	})

	server.OnExecute(0xFFAA, func(client *giotgo.ClientHandler, data giot_packet.Data) *bytes.Buffer {
		fmt.Println(data)
		return nil
	})

	server.ClientAuth(func(username, password string) bool {
		fmt.Println(username, password)
		return true
	})

	server.Serve("localhost:2000")
}
