package main

import (
	"fmt"

	giotgo "github.com/ghuvrons/g-IoT-Go"
)

func main() {
	// var crc uint32 = 0
	// crc = crc32.Update(crc, crc32.IEEETable, []byte("ghuvrons"))

	// fmt.Printf("0x%.8X\r\n", crc)

	// return
	fmt.Println("Start")

	var server = giotgo.NewServer()

	setCmdHandlers(server)

	server.ClientAuth(func(username, password string) bool {
		fmt.Println(username, password)
		return true
	})

	server.Serve("localhost:2000")
}
