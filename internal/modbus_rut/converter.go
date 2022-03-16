package modbus_rut

import (
	"encoding/binary"
	"encoding/hex"
	"math"
)

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

func modbusConvertService(results []byte, registerMap RegMap) []string {
	return []string{"1", "2"}
}
