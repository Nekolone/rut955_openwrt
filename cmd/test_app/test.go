package main

import (
	"encoding/binary"
	"encoding/hex"
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
	log.Println(getFloat64(res, 1, 2, 3, 4, 5, 6, 7, 8))
	log.Println(getFloat64(res, 7, 8, 5, 6, 3, 4, 1, 2))
}

func getSigInt8(bytes []byte) []int8 {
	var resul []int8
	for _, d := range bytes {
		resul = append(resul, int8(d))
	}
	return resul
}

func getUnsInt8(bytes []byte) []uint8 {
	return bytes
}

func getSigInt16(bytes []byte, p1 int, p2 int) []int16 {
	var resul []int16
	for i := 1; i <= len(bytes)-2; i += 2 {
		resul = append(resul, byteToInt16([]byte{bytes[i-(2-p1)], bytes[i-(2-p2)]}))
	}
	return resul
}

func getUnsInt16(bytes []byte, p1 int, p2 int) []uint16 {
	var resul []uint16
	for i := 1; i <= len(bytes)-2; i += 2 {
		resul = append(resul, byteToUint16([]byte{bytes[i-(2-p1)], bytes[i-(2-p2)]}))
	}
	return resul
}

func getSigInt32(bytes []byte, p1 int, p2 int, p3 int, p4 int) []int32 {
	var resul []int32
	for i := 1; i <= len(bytes)-2; i += 2 {
		resul = append(resul, byteToInt32([]byte{bytes[i-(4-p1)], bytes[i-(4-p2)], bytes[i-(4-p3)], bytes[i-(4-p4)]}))
	}
	return resul
}

func getUnsInt32(bytes []byte, p1 int, p2 int, p3 int, p4 int) []uint32 {
	var resul []uint32
	for i := 1; i <= len(bytes)-2; i += 2 {
		resul = append(resul, byteToUint32([]byte{bytes[i-(4-p1)], bytes[i-(4-p2)], bytes[i-(4-p3)], bytes[i-(4-p4)]}))
	}
	return resul
}

func getFloat32(bytes []byte, p1 int, p2 int, p3 int, p4 int) []float32 {
	var resul []float32
	for i := 3; i <= len(bytes)-4; i += 4 {
		resul = append(resul, byteToFloat32([]byte{bytes[i-(4-p1)], bytes[i-(4-p2)], bytes[i-(4-p3)], bytes[i-(4-p4)]}))
	}
	return resul
}

func getFloat64(bytes []byte, p1 int, p2 int, p3 int, p4 int, p5 int, p6 int, p7 int, p8 int) []float64 {
	var resul []float64
	for i := 3; i <= len(bytes)-4; i += 4 {
		resul = append(resul, byteToFloat64([]byte{bytes[i-(8-p1)], bytes[i-(8-p2)], bytes[i-(8-p3)], bytes[i-(8-p4)],
			bytes[i-(8-p5)], bytes[i-(8-p6)], bytes[i-(8-p7)], bytes[i-(8-p8)]}))
	}
	return resul
}

func getAscii(bytes []byte) string {
	return byteToAscii(bytes)
}

func getHex(bytes []byte) []string {
	var resul []string
	for i := 1; i < len(bytes)-2; i += 2 {
		resul = append(resul, byteToHex([]byte{bytes[i-1], bytes[i]}))
	}
	return resul
}

func byteToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func byteToAscii(bytes []byte) string {
	return string(bytes)
}

func byteToInt16(bytes []byte) int16 {
	return int16(binary.BigEndian.Uint16(bytes))
}

func byteToUint16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

func byteToInt32(bytes []byte) int32 {
	return int32(binary.BigEndian.Uint32(bytes))
}
func byteToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func byteToFloat32(bytes []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(bytes))
}

func byteToFloat64(bytes []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(bytes))
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
