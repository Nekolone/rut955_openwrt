package main

import (
	"github.com/goburrow/modbus"
	"log"
)

func main() {
	// Modbus TCP
	client := modbus.TCPClient("172.16.1.12:502")
	// Read input register 9
	results, err := client.ReadHoldingRegisters(12, 20)

	if err != nil {
		log.Println(err)
	}

	log.Println(results)

}
