package wialonClient

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func sendData(data string, clientConnection *net.Conn, networkStatus *string, bufferPath string) {
	switch *networkStatus {
	case "buffering":
		saveToBuffer(data, bufferPath)
	case "online":
		if send(data, clientConnection, networkStatus) != "success" {
			saveToBuffer(data, bufferPath)
			return
		}
	default:
		log.Printf("unexpected error, sendData postBUffering stage (impossible) %s", *networkStatus)
	}
}

func send(data string, clientConnection *net.Conn, networkStatus *string) (answer string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Send panic. Recover msg > %v", r)
			if *clientConnection != nil {
				_ = (*clientConnection).Close()
			}
			log.Print("networkStatus -> buffering")
			*networkStatus = "buffering"
			answer = fmt.Sprint(r)
		}
	}()
	serverReader := bufio.NewReader(*clientConnection)
	timer := time.NewTimer(time.Second * 90)
	result := make(chan string)
	defer func() {
		if timer != nil {
			timer.Stop()
		}
		close(result)
	}()
	log.Print("Send to server")
	for i := 1; i < 5; i++ {
		if _, err := (*clientConnection).Write([]byte(data + "\r\n")); err != nil {
			if *clientConnection != nil {
				_ = (*clientConnection).Close()
			}
			log.Panicf("Write to server failed: %v", err.Error())
		}

		timer = time.NewTimer(time.Second * 20)

		go func() {
			serverResponse, err := serverReader.ReadString('\n')
			if err != nil {
				result <- fmt.Sprintf("msg:%s error: %v\n", serverResponse, err)
			}
			result <- serverResponse
		}()

		select {
		case res := <-result:
			switch reactToResponse(res) {
			case "resend":
				log.Print("resend data")
			case "success":
				log.Print("data successfully sent to server")
				return "success"
			case "#RESTART":
				*networkStatus = "RESTART"
				log.Print("networkStatus -> RESTART")
				return "success"
			case "#STOP":
				*networkStatus = "RESTART"
				log.Print("networkStatus -> RESTART")
				return "success"
			default:
				log.Printf("unknown answer %s", res)
			}
		case <-timer.C:
			log.Print("timeout")
		}
	}
	if (*clientConnection) != nil {
		_ = (*clientConnection).Close()
	}
	log.Panicf("send error")
	return
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
	case "#RESTART\r\n":
		return "RESTART"
	case "#STOP\r\n":
		return "RESTART"
	default:
		return response
	}
}
