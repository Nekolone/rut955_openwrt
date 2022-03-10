package main

import (
	"encoding/binary"
	"github.com/goburrow/modbus"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// Modbus TCP
	res := getModbusData()
	log.Println(res[0:4])
	for i := 8; i <= len(res)-8; i += 8 {
		log.Println(i)
		resul := Float64frombytes(res[i-8 : i])
		log.Println(resul)
	}
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func getModbusData() []byte {
	handler := modbus.NewTCPClientHandler("172.16.1.12:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0x0C
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	if err != nil {
		log.Println(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(0, 20)
	return results
}
