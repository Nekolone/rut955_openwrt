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
	//tcpAddr, err := net.ResolveTCPAddr("tcp", "123.0.0.1:9000")
	//clientConnection, _ := net.DialTCP("tcp", nil, tcpAddr)
	var nilCOn net.Conn
	nilCOn = nil
	clientConnection, err := net.Dial("tcp", "127.0.0.1:9000")
	log.Print(err)
	log.Print(&err)
	log.Print(clientConnection)

	log.Println("test")
	log.Print(testF(&nilCOn))
	log.Println("succ")

}

func testF(clientConnection *net.Conn) *net.Conn {
	defer func() {
		if r := recover(); r!=nil{

		}
	}()

	log.Print(1)
	log.Panicf("HIHIHAHA")
	return clientConnection
}
