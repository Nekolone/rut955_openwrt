package wialonClient

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Config struct {
	WialonServerAddress string `json:"wialon_server_address"`
	ConnectionType      string `json:"connection_type"`
	DataBufferPath      string `json:"data_buffer_path"`
	Login               string `json:"login"`
	Password            string `json:"password"`
}

func Start(dataChan chan string, conf *Config) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("recover in wialon client. Panic > %v", r)
			if strings.Contains(fmt.Sprintf("%v", r), "FATAL") {
				log.Panicf("FATAL error in data processing service. Use painc. Reason: %v", r)
			}
		}
	}()
	log.Print("Wialon Client start")
	networkStatus := "start" // var that describes the state of the connection
	clientConnection, addr := ConnectToServer(conf, &networkStatus)
	defer func() {
		if *clientConnection != nil {
			_ = (*clientConnection).Close()
		}
	}()

	log.Print("Wialon Client routines start")
	done := make(chan string)
	go ReconnectingService(conf, addr, clientConnection, &networkStatus, done)
	go DataWorker(conf, clientConnection, &networkStatus, dataChan, done)

	log.Print("Wialon Client wait for routines")
	if d := <-done; true {
		go func() {
			for networkStatus != "DONE" {
				networkStatus = "RESTART"
				time.Sleep(100 * time.Millisecond)
			}
		}()
		log.Printf("first done triggered. Restart wialon client. Reason : %v", d)
		d = <-done
		networkStatus = "DONE"
		close(done)
		log.Panicf("Restart wialon client. Reason: %v", d)
	}
}

func ConnectToServer(conf *Config, networkStatus *string) (con *net.Conn, addr string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("First connection error. Save all data into buffer file. err msg > %v", r)
			if strings.Contains(fmt.Sprintf("%v", r), "FATAL") {
				log.Panicf("FATAL error in data processing service. Use painc. Reason: %v", r)
			}
			*networkStatus = "buffering"
			log.Print("networkStatus -> buffering")
		}
	}()
	addr = conf.WialonServerAddress
	var nilCon net.Conn
	nilCon = nil
	//log.Printf("resolving tcp addr %v", conf.WialonServerAddress)
	//if tcpAddr, err = net.ResolveTCPAddr(conf.ConnectionType, conf.WialonServerAddress); err != nil {
	//	log.Panicf("FATAL : cant resolve tcp addr > %v", err)
	//}
	//log.Print("resolving successfully")

	log.Print("connecting to wialon server")
	clientConnection, err := net.Dial(conf.ConnectionType, conf.WialonServerAddress)
	if err != nil {
		con = &nilCon
		log.Panicf("dial error %v", err)
	}
	log.Print("connecting successfully")

	log.Print("login to wialon server")
	if answer := login(&clientConnection, conf.Login, conf.Password); answer != "" {
		con = &nilCon
		log.Panicf(answer)
	}
	con = &clientConnection
	log.Print("login successfully")

	log.Print("check buffer file and start data sharing")
	*networkStatus = "postBuffering"
	log.Print("networkStatus -> postBuffering")
	return
}

func DataWorker(conf *Config, clientConnection *net.Conn, networkStatus *string, dataChan chan string, done chan string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in f > %v", r)
		}
		done <- "module restart"
		*networkStatus = "RESTART"

	}()
	log.Println("DataWorker start")

	for {
		data := <-dataChan
		log.Print(*networkStatus)
		log.Print(clientConnection)
		log.Print(*clientConnection)
		log.Print(data)
		switch *networkStatus {
		case "online":
			sendData(data, clientConnection, networkStatus, conf.DataBufferPath)
		case "buffering":
			saveToBuffer(data, conf.DataBufferPath)
		case "postBuffering":
			sendBufferData(clientConnection, networkStatus, conf.DataBufferPath)
			sendData(data, clientConnection, networkStatus, conf.DataBufferPath)
		case "stop":
			log.Println("client service stop")
			return
		case "RESTART":
			return
		default:
			saveToBuffer(data, conf.DataBufferPath)
		}
	}
}
