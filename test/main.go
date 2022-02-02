package main

import (
	"fmt"

	giotgo "github.com/ghuvrons/g-IoT-Go"
)

func main() {
	fmt.Println("Start")
	var server = giotgo.Server{}
	server.Serve("localhost:2000")
}
