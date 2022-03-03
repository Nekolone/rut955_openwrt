package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	data := make(chan int, 1000)
	go hehe(data)
	time.Sleep(20 * time.Second)
	dd := []int{1}
	for len(dd) < 10000000 {
		dd = append(dd, <-data)
	}
	log.Println(dd)
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
