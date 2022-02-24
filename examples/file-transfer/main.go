package main

import (
	"fmt"

	giotgo "github.com/ghuvrons/g-IoT-Go"
)

func main() {
	fmt.Println("Start")

	var server = giotgo.NewServer()

	setCmdHandlers(server)

	server.ClientAuth(func(username, password string) bool {
		fmt.Println(username, password)
		return true
	})

	server.Serve("localhost:2000")
}
