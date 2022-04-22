package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	log.Println("test")
	log.Print(testF())
	log.Println("succ")

}

func testF() (answer string) {
	var clientConnection *net.TCPConn
	defer func() {
		if r := recover(); r != nil {
			log.Print("networkStatus -> buffering")
			answer = fmt.Sprint(r)
			if clientConnection != nil {
				clientConnection.Close()
			}

		}
	}()

	log.Print(1)
	log.Panicf("panicEND")
	return
}
