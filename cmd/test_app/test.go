package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	list := []string{"4", "3", "2", "1"}
	log.Println(len(list))
	log.Println(list[len(list)-1])
}

func hehe(data chan int) {
	for i := 1; i < 10000000; i++ {
		fmt.Println(i)
		data <- i
	}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("cant close connection")
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
