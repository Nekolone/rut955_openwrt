package wialonClient

import (
	"bufio"
	"log"
	"net"
	"os"
)

func sendBufferData(clientConnection *net.TCPConn, networkStatus *string, bufferPath string) {
	log.Print("Send data from buffer file")
	fileHandler, err := os.OpenFile(bufferPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("no buffer file")
		*networkStatus = "online"
		log.Println("networkStatus -> online")
		return
	}

	newBufferPath := "/tmp/RWG_app_buffer/newBuffer.buf"
	defer deleteOldBuffer(newBufferPath, bufferPath)
	defer fileHandler.Close()

	fileScanner := bufio.NewScanner(fileHandler)
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
		log.Panicf("scan file error: %v", err)
		return
	}

	if *networkStatus == "postBuffering" {
		log.Println("buffered data send successfully")
		*networkStatus = "online"
		log.Println("networkStatus -> online")
		return
	}
	log.Println("still buffering")
}

func saveToBuffer(data string, bufferPath string) {
	fileHandler, err := os.OpenFile(bufferPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("open error %v, try to create %s\n", err, bufferPath)
		if fileHandler, err = os.Create(bufferPath); err != nil {
			log.Println("create error", err)
			os.Exit(1)
		}
	}
	defer fileHandler.Close()
	log.Print("Save to buffer")
	_, _ = fileHandler.WriteString(data + "\n")
}

func deleteOldBuffer(newBufferPath string, bufferPath string) {
	err := os.Remove(bufferPath)
	if err != nil {
		log.Printf("remove buffer error> %v", err)
	}
	err = os.Rename(newBufferPath, bufferPath)
	if err != nil {
		log.Printf("rename buffer error> %v", err)
	}
}
