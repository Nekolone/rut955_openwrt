package custom

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func listenService(serverConnection net.Listener, deviceDataChan chan map[string][]string, serviceName string) {
	for {
		deviceConnection, err := serverConnection.Accept()
		if err != nil {
			log.Println("listen service error")
		}
		go handleRequest(deviceConnection, deviceDataChan, serviceName)
	}
}

func getCurTime() string {
	out := []byte(time.Now().Format("2006-01-02 15:04:05.000000"))
	return string(out[8:10]) + string(out[5:7]) + string(out[2:4]) + string(out[11:13]) + string(out[14:16]) + string(out[17:19]) + string(out[20:26])
}

func handleRequest(connection net.Conn, deviceDataChan chan map[string][]string, serviceName string) {
	defer connection.Close()
	clientReader := bufio.NewReader(connection)
	for {
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest = strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			}
			deviceDataChan <- map[string][]string{
				serviceName: {
					getCurTime(),
					fmt.Sprint(clientRequest),
				},
			}
			log.Println(clientRequest)

		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
	}
}
