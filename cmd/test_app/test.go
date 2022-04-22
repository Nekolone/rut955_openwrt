package main

import (
	"log"
	"net"
)

func main() {
	serverConnection, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Panicf("listenService error msg> %v", err)
	}
	defer serverConnection.Close()

	go func() {
		_, _ = serverConnection.Accept()
	}()
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9000")
	clientConnection, _ := net.DialTCP("tcp", nil, tcpAddr)
	// clientConnection, err := net.Dial("tcp", "127.0.0.1:9000")
	log.Print(err)

	log.Println("test")
	log.Print(testF(&clientConnection))
	log.Println("succ")

}

func testF(clientConnection **net.TCPConn) (answer string) {
	defer func() {
		_ = (*clientConnection).Close().Error()
	}()

	log.Print(1)
	// log.Panicf("panicEND")
	return
}
