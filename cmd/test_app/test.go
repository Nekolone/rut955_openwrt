package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	modbusclient "rut955_openwrt/pkg/go-modbus"
)

func main() {
	var (
		host string
		port int
	)
	const (
		defaultHost = "172.16.1.12"
		defaultPort = modbusclient.MODBUS_PORT
	)

	flag.StringVar(&host, "host", defaultHost, "Slave device host (url or ip address)")
	flag.IntVar(&port, "port", defaultPort, fmt.Sprintf("Slave device port (the default is %d)", defaultPort))
	flag.Parse()
	trace := true

	conn, cerr := modbusclient.ConnectTCP(host, port)
	if cerr != nil {
		log.Println(fmt.Sprintf("Connection error: %s", cerr))
		os.Exit(1)
	}
	// attempt to read one (0x01) holding registers starting at address 200
	readData := make([]byte, 3)
	readData[0] = byte(200 >> 8)   // (High Byte)
	readData[1] = byte(200 & 0xff) // (Low Byte)
	readData[2] = 12

	// make this read request transaction id 1, with a 300 millisecond tcp timeout
	readResult, readErr := modbusclient.TCPRead(conn, 3000, 1, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x0C, readData, trace)
	if readErr != nil {
		log.Println(readErr)
	}
	log.Println(readResult)

	modbusclient.DisconnectTCP(conn)

}

func connectModbus() {

}
