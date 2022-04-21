package wialonClient

import (
	"log"
	net "net"
	"time"
)

func ReconnectingService(conf *Config, tcpAddr **net.TCPAddr, clientConnection **net.TCPConn, networkStatus *string, done chan string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in f > %v", r)
		}
		done <- "module restart"
		*networkStatus = "RESTART"
		(*clientConnection).Close()
	}()
	recTimer := time.NewTicker(time.Minute)
	log.Println("ReconnectingService start")
	var err error
	for range recTimer.C {
		if *networkStatus == "RESTART"{
			return
		}
		if *networkStatus != "buffering" {
			continue
		}
		log.Print("reconnecting to wialon server")
		*clientConnection, err = net.DialTCP(conf.ConnectionType, nil, *tcpAddr)
		if err != nil {
			log.Printf("Reconnecting failed: %v", err.Error())
			continue
		}
		log.Print("reconnect successfully")

		log.Print("login to wialon server")
		if res := login(*clientConnection, conf.Login, conf.Password); res != "" {
			log.Println("login error")
			_ = (*clientConnection).Close()
			*networkStatus = "buffering"
			log.Println("networkStatus -> buffering")
			continue
		}
		log.Print("login successfully")

		*networkStatus = "postBuffering"
		log.Println("networkStatus -> postBuffering")
	}
}
