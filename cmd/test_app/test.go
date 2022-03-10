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
	log.Println(res)
	log.Println(string(res))
	for i := 0; i <= len(res)-4; i += 4 {
		resul := Float64frombytes(res[i : i+4])
		log.Println(resul)
	}
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
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
	results, err := client.ReadHoldingRegisters(1, 10)
	return results
}
