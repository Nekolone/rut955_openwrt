package client

import (
	"log"
	"net"
)

func ConnectToServer(servAddr string, network string, networkStatus *string, id string, pass string) (
	*net.TCPConn, *net.TCPAddr) {

	log.Printf("connecting to server %v", servAddr)

	tcpAddr, err := net.ResolveTCPAddr(network, servAddr)
	if err != nil {
		log.Fatal("Client creation error")
	}

	clientConnection, err := net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		log.Println("Dial failed:", err.Error())
		*networkStatus = "buffering"
		log.Println("networkStatus -> buffering")
		return clientConnection, tcpAddr
	}

	res := login(&clientConnection, id, pass)
	if res != "" {
		log.Printf("login error: %s\n", res)
		clientConnection.Close()
		*networkStatus = "buffering"
		log.Println("networkStatus -> buffering")
		return clientConnection, tcpAddr
	}

	log.Printf("connecting successfully")
	*networkStatus = "postBuffering"
	log.Println("networkStatus -> postBuffering")
	return clientConnection, tcpAddr
}

func DataWorker(networkStatus *string, clientConnection **net.TCPConn, dataChan chan string, bufferPath string) {
	log.Println("DataWorker start")

	for {
		data := <-dataChan
		switch *networkStatus {
		case "online":
			sendData(data, *clientConnection, networkStatus, bufferPath)
		case "buffering":
			saveToBuffer(data, bufferPath)
		case "postBuffering":
			sendBufferData(*clientConnection, networkStatus, bufferPath)
			sendData(data, *clientConnection, networkStatus, bufferPath)
		case "stop":
			log.Println("client service stop")
			return
		case "restart":
			log.Println("client service restart")
			return
		default:
			saveToBuffer(data, bufferPath)
		}
	}
}
