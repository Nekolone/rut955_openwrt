package wialonClient

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func login(clientConnection **net.TCPConn, id string, pass string) (answer string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recover in f > %v", r)
			answer = fmt.Sprint(r)
		}
	}()
	serverReader := bufio.NewReader(*clientConnection)
	var timer *time.Timer
	result := make(chan string)
	defer func() {
		if timer != nil {
			timer.Stop()
		}
		close(result)
	}()

	for i := 0; i < 5; i++ {
		if _, err := (*clientConnection).Write([]byte(string("#L#" + id + ";" + pass + "\r\n"))); err != nil {
			if *clientConnection != nil {
				_ = (*clientConnection).Close()
			}
			log.Panicf("close connection\nlogin error %s", err)
		}

		timer = time.NewTimer(time.Second * 30)

		go func() {
			serverResponse, err := serverReader.ReadString('\n')
			if err != nil {
				result <- fmt.Sprintf("msg:%s error: %v\n", serverResponse, err)
			}
			result <- serverResponse
		}()

		select {
		case res := <-result:
			switch reactToLoginResponse(res) {
			case "resend":
				log.Printf("login -%d- failed", i)
			case "success":
				return ""
			default:
				if *clientConnection != nil {
					_ = (*clientConnection).Close()
				}
				log.Panicf("close connection\nlogin error. Answer %s", res)
			}
		case <-timer.C:
			log.Print("timeout")
		}
	}
	if *clientConnection != nil {
		_ = (*clientConnection).Close()
	}
	log.Panicf("close connection\nlogin error. Wrong login or pass:%v %v", id, pass)
	return
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
