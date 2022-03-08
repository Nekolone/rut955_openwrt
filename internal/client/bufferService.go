package client

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/exec"
)

func sendBufferData(clientConnection *net.TCPConn, networkStatus *string, bufferPath string) {
	log.Println("send buffered data")

	fileHanler, err := os.OpenFile("buffer.buf", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("cant read buffer file")
		*networkStatus = "online"
		log.Println("networkStatus -> online")
		return
	}

	newBufferPath := "newBuffer.buf"
	defer deleteOldBuffer(newBufferPath, bufferPath)
	defer fileHanler.Close()

	fileScanner := bufio.NewScanner(fileHanler)
	var msg string
	for fileScanner.Scan() {
		msg = fileScanner.Text()
		switch *networkStatus {
		case "postBuffering":
			if send(msg, clientConnection, networkStatus) != "success" {
				saveToBuffer(msg, newBufferPath)
			}
		case "buffering":
			saveToBuffer(msg, newBufferPath)
		case "online":
			if send(msg, clientConnection, networkStatus) != "success" {
				saveToBuffer(msg, newBufferPath)
			}
		default:
			log.Printf("buffer unexpected error, networkStatus is not postBuffering or buffering > %s", *networkStatus)
			return
		}
	}

	if err = fileScanner.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}

	if *networkStatus == "postBuffering" {
		log.Println("buffered data send successfully")
		*networkStatus = "online"
		log.Println("networkStatus -> online")
	}
	log.Println("still buffering")
}

func saveToBuffer(data string, bufferPath string) {
	fileHandler, err := os.OpenFile(bufferPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("open error %v, try to create buffer.buf\n", err)
		if fileHandler, err = os.Create("buffer.buf"); err != nil {
			log.Println("create error", err)
			os.Exit(1)
		}
		exec.Command("chmod +x buffer.buf")
	}
	defer fileHandler.Close()
	fileHandler.WriteString(data + "\n")
}

func deleteOldBuffer(bufferPath string, newBufferPath string) {
	err := os.Remove(bufferPath)
	if err != nil {
		log.Printf("remove buffer error> %v", err)
	}
	err = os.Rename(newBufferPath, bufferPath)
	if err != nil {
		log.Printf("rename buffer error> %v", err)
	}
}
