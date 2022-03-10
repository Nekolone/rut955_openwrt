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
	log.Println(res[0:8])
	log.Println(getSigInt8b(res))
}

func getUnsInt8b(bytes []byte) []uint8 {
	return bytes
}

func getSigInt8b(bytes []byte) []int8 {
	var resul []int8
	for _, d := range bytes {
		log.Println(int8(d))
		resul = append(resul, int8(d))
	}
	return resul

}

func getFloat32bytes3412(bytes []byte, p1 int, p2 int, p3 int, p4 int) []float32 {
	var resul []float32
	for i := 3; i <= len(bytes)-4; i += 4 {
		resul = append(resul, Float32frombytes([]byte{bytes[i-(4-p1)], bytes[i-(4-p2)], bytes[i-(4-p3)], bytes[i-(4-p4)]}))
	}
	return resul
}

func Float32frombytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
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
