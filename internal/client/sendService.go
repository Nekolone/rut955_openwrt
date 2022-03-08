package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func sendData(data string, clientConnection *net.TCPConn, networkStatus *string, bufferPath string) {
	switch *networkStatus {
	case "buffering":
		saveToBuffer(data, bufferPath)
	case "online":
		if send(data, clientConnection, networkStatus) != "success" {
			saveToBuffer(data, bufferPath)
		}
	default:
		log.Printf("unexpected error, sendData postBUffering stage (impossible) %s", *networkStatus)
	}
}

func send(data string, clientConnection *net.TCPConn, networkStatus *string) string {
	serverReader := bufio.NewReader(clientConnection)
	timer := time.NewTimer(time.Second * 90)
	for i := 1; i < 5; i++ {
		timer = time.NewTimer(time.Second * 90)

		if _, err := clientConnection.Write([]byte(data + "\r\n")); err != nil {
			log.Println("Write to server failed:", err.Error())
			*networkStatus = "buffering"
			log.Println("networkStatus -> buffering")
			return "send error"
		}

		result := make(chan string)

		go func() {
			serverResponse, err := serverReader.ReadString('\n')
			switch err {
			case nil:
				result <- serverResponse
			case io.EOF:
				result <- "server closed the connection"
			default:
				result <- fmt.Sprintf("server error: %v", err)
			}
		}()

		select {
		case res := <-result:
			switch reactToResponse(res) {
			case "resend":
				log.Println("resend data")
			case "success":
				log.Println("data successfully sent to server")
				return "success"
			case "#RESTART":
				*networkStatus = "restart"
				log.Println("networkStatus -> restart")
				return "success"
			case "#STOP":
				*networkStatus = "stop"
				log.Println("networkStatus -> stop")
				return "success"
			default:
				log.Printf("unknown answer %s", res)
			}
		case <-timer.C:
			log.Println("timeout")
		}
	}
	timer.Stop()

	log.Println("send error")
	clientConnection.Close()
	*networkStatus = "buffering"
	log.Println("networkStatus -> buffering")
	return "send error"
}

func reactToResponse(response string) string {
	switch response {
	case "#AD#-1\r\n":
		return "resend"
	case "#AD#0\r\n":
		return "success"
	case "#AD#1\r\n":
		return "success"
	case "#AD#10\r\n":
		return "success"
	case "#AD#11\r\n":
		return "success"
	case "#AD#12\r\n":
		return "success"
	case "#AD#13\r\n":
		return "success"
	case "#AD#14\r\n":
		return "success"
	case "#AD#15\r\n":
		return "success"
	default:
		return response
	}
}
