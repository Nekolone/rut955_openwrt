package dataProcessingService

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

func listenService(serverConnection net.Listener, deviceDataChan chan string) {
	for {
		deviceConnection, err := serverConnection.Accept()
		if err != nil {
			log.Println("listen service error")
		}
		go handleRequest(deviceConnection, deviceDataChan)
	}
}

func handleRequest(connection net.Conn, deviceDataChan chan string) {
	defer connection.Close()
	clientReader := bufio.NewReader(connection)
	for {
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest = strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			}
			deviceDataChan <- clientRequest
			log.Println(clientRequest)

		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
	}
}
