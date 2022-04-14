package wialonClient

import (
	"log"
	"net"
	"time"
)

func ReconnectingService(conf *Config, tcpAddr **net.TCPAddr, clientConnection **net.TCPConn, networkStatus *string) {
	recTimer := time.NewTicker(time.Minute)
	log.Println("ReconnectingService start")
	var err error
	for range recTimer.C {
		if *networkStatus != "buffering" {
			continue
		}
		log.Println("reconnecting")
		*clientConnection, err = net.DialTCP(conf.ConnectionType, nil, *tcpAddr)
		if err != nil {
			log.Println("Reconnecting failed: ", err.Error())
			continue
		}
		res := login(clientConnection, conf.Login, conf.Password)
		if res != "" {
			log.Println("login error")
			_ = (*clientConnection).Close()
			*networkStatus = "buffering"
			log.Println("networkStatus -> buffering")
			continue
		}
		log.Println("reconnect successfully")
		*networkStatus = "postBuffering"
		log.Println("networkStatus -> postBuffering")
	}
}
