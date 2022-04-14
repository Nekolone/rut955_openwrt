package wialonClient

import (
	"log"
	"net"
)

type Config struct {
	WialonServerAddress string `json:"wialon_server_address"`
	ConnectionType      string `json:"connection_type"`
	DataBufferPath      string `json:"data_buffer_path"`
	Login               string `json:"login"`
	Password            string `json:"password"`
}

func Start(dataChan chan string, conf *Config) {
	networkStatus := "start"

	log.Println("wialon cli - start")

	clientConnection, tcpAddr := ConnectToServer(conf, &networkStatus)
	go ReconnectingService(conf, &tcpAddr, &clientConnection, &networkStatus)
	DataWorker(conf, &networkStatus, &clientConnection, dataChan)

	log.Println("wialon cli - end")
}

func ConnectToServer(conf *Config, networkStatus *string) (*net.TCPConn, *net.TCPAddr) {
	log.Printf("wialon cli - connecting to server %v", conf.WialonServerAddress)

	tcpAddr, err := net.ResolveTCPAddr(conf.ConnectionType, conf.WialonServerAddress)
	if err != nil {
		log.Fatal("Client creation error")
	}

	clientConnection, err := net.DialTCP(conf.ConnectionType, nil, tcpAddr)
	if err != nil {
		log.Println("Dial failed:", err.Error())
		*networkStatus = "buffering"
		log.Println("networkStatus -> buffering")
		return clientConnection, tcpAddr
	}

	log.Println("wialon cli - login")

	res := login(&clientConnection, conf.Login, conf.Password)
	if res != "" {
		log.Printf("login error: %s\n", res)
		_ = clientConnection.Close()
		*networkStatus = "buffering"
		log.Println("networkStatus -> buffering")
		return clientConnection, tcpAddr
	}

	log.Printf("connecting successfully")
	*networkStatus = "postBuffering"
	log.Println("networkStatus -> postBuffering")
	return clientConnection, tcpAddr
}

func DataWorker(conf *Config, networkStatus *string, clientConnection **net.TCPConn, dataChan chan string) {
	log.Println("DataWorker start")

	for {
		data := <-dataChan
		switch *networkStatus {
		case "online":
			sendData(data, *clientConnection, networkStatus, conf.DataBufferPath)
		case "buffering":
			saveToBuffer(data, conf.DataBufferPath)
		case "postBuffering":
			sendBufferData(*clientConnection, networkStatus, conf.DataBufferPath)
			if r := recover(); r != nil {
				log.Printf("Recovered in f > %v", r)
			}
			sendData(data, *clientConnection, networkStatus, conf.DataBufferPath)
		case "stop":
			log.Println("client service stop")
			return
		case "restart":
			log.Println("client service restart")
			return
		default:
			saveToBuffer(data, conf.DataBufferPath)
		}
	}
}
