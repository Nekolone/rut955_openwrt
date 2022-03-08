package client

import (
	"log"
	"net"
	"time"
)

func ReconnectingService(tcpAddr **net.TCPAddr, network string, clientConnection **net.TCPConn, networkStatus *string,
	id string, pass string) {
	recTimer := time.NewTicker(time.Minute)
	log.Println("ReconnectingService start")
	var err error
	for range recTimer.C {
		if *networkStatus == "buffering" {
			log.Println("reconnecting")
			*clientConnection, err = net.DialTCP(network, nil, *tcpAddr)
			if err != nil {
				log.Println("Reconnecting failed: ", err.Error())
				continue
			}
			res := login(clientConnection, id, pass)
			if res != "" {
				log.Println("login error")
				(*clientConnection).Close()
				*networkStatus = "buffering"
				log.Println("networkStatus -> buffering")
				continue
			}
			log.Println("reconnect successfully")
			*networkStatus = "postBuffering"
			log.Println("networkStatus -> postBuffering")
		}
	}
}
