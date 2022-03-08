package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func login(clientConnection **net.TCPConn, id string, pass string) string {
	serverReader := bufio.NewReader(*clientConnection)
	var timer *time.Timer
	for i := 1; i < 5; i++ {
		timer = time.NewTimer(time.Second * 30)

		if _, err := (*clientConnection).Write([]byte(string("#L#" + id + ";" + pass + "\r\n"))); err != nil {
			log.Println("Login connection error:", err.Error())
			return "error"
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
				result <- fmt.Sprintf("server error: %v\n", err)
			}
		}()

		select {
		case res := <-result:
			switch reactToLoginResponse(res) {
			case "resend":
				log.Println("ligin problem")
				break
			case "success":
				log.Println("login successfully")
				return ""
			default:
				return "error"
			}
		case <-timer.C:
			log.Println("timeout")
		}
	}
	timer.Stop()
	return "error"
}

func reactToLoginResponse(response string) string {
	switch response {
	case "#AL#0\r\n":
		return "resend"
	case "#AL#1\r\n":
		return "success"
	default:
		return response
	}
}
