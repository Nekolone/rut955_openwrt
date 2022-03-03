package client

import (
	"log"
	"net"
	"time"
)

func Start(servAddr string, network string, networkStatus *string) (*net.TCPConn, *net.TCPAddr) {
	tcpAddr, err := net.ResolveTCPAddr(network, servAddr)
	if err != nil {
		log.Fatal("Client creation error")
	}

	clientConnection, err := net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		*networkStatus = "buffering"
		log.Println("Dial failed:", err.Error())
		return clientConnection, tcpAddr
	}
	*networkStatus = "online"
	return clientConnection, tcpAddr
}

func ReconnectingService(tcpAddr **net.TCPAddr, network string, clientConnection **net.TCPConn, networkStatus *string) {
	recTimer := time.NewTicker(15 * time.Second)
	log.Println("ReconnectingService start")
	var err error
	for range recTimer.C {
		if *networkStatus == "buffering" {
			*clientConnection, err = net.DialTCP(network, nil, *tcpAddr)
			if err != nil {
				log.Println("Reconnecting failed:", err.Error())
				continue
			}
			*networkStatus = "online"
			log.Println("reconnect successfully")
		}
	}
}

func DataWorker(networkStatus *string, clientConnection **net.TCPConn, dataChan chan string) {
	log.Println("DataWorker Start")
	for {
		log.Println("wait 4  data, status>", *networkStatus)
		data := <-dataChan
		log.Println("no data")
		switch *networkStatus {
		case "online":
			sendData(data, *clientConnection, networkStatus)
		case "buffering":
			saveToBuffer(data)
		case "postBuffering":
			sendBufferData(data, *clientConnection, networkStatus)
			sendData(data, *clientConnection, networkStatus)
		case "Stop":
			log.Println("client service stop")
			return
		//case "Restart":
		//	log.Println("client service restart")
		//	return
		default:
			saveToBuffer(data)
		}
	}
}

func saveToBuffer(data string) {
	//сохранить в буфеп
}

func sendBufferData(data string, clientConnection *net.TCPConn, networkStatus *string) {
	//берет из файла данные и вызывает функцию send
	//for EOF{send(data, clientConnection, networkStatus)}
	*networkStatus = "online"
}

func sendData(data string, clientConnection *net.TCPConn, networkStatus *string) {
	send(data, clientConnection, networkStatus)
}

func send(data string, clientConnection *net.TCPConn, networkStatus *string) {
	status := "start"
	for i := 1; i < 5; i++ {
		log.Println([]byte(data))
		_, err := clientConnection.Write([]byte(data))
		if err != nil {
			log.Println("Write to server failed:", err.Error())
			status = "Error"
			saveToBuffer(data)
			break
		}
		log.Println("send> ", data)
		reply := make([]byte, 1024)
		_, err = clientConnection.Read(reply)

		if err != nil {
			log.Println("Client Read error: ", err)
			status = "Error"
			break
		}

		if err == nil {
			status = "Successfully"
			switch string(reply) {
			case "1":
				//отправить повторно
				break
			case "2":
				//отправить повторно, заменив часть
				break
			case "3":
				log.Println("успех")
				//УСПЕХ
				break
			case "4":
				break
			case "5":
				break
			case "#RESTART":
				*networkStatus = "Restart"
			case "#STOP":
				*networkStatus = "Stop"
			default:
				log.Println("wrong answer", string(reply))

			}
		}
	}
	if status == "Error" {
		log.Println("error -> set new status buff")
		*networkStatus = "buffering"
	}
}
