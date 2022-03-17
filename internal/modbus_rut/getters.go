package modbus_rut

import "log"

func GetSigInt8(d byte) (result int8) {
	result = int8(d)
	return result
}

func GetUnsInt8(bytes byte) (result uint8) {
	return bytes
}

func GetSigInt16(bytes []byte, p1 int, p2 int) (result int16) {
	result = byteToInt16([]byte{bytes[p1-1], bytes[p2-1]})
	return result
}

func GetUnsInt16(bytes []byte, p1 int, p2 int) (result uint16) {
	result = byteToUint16([]byte{bytes[p1-1], bytes[p2-1]})
	return result
}

func GetSigInt32(bytes []byte, p1 int, p2 int, p3 int, p4 int) (result int32) {
	result = byteToInt32([]byte{bytes[p1-1], bytes[p2-1], bytes[p3-1], bytes[p4-1]})
	return result
}

func GetUnsInt32(bytes []byte, p1 int, p2 int, p3 int, p4 int) (result uint32) {
	result = byteToUint32([]byte{bytes[p1-1], bytes[p2-1], bytes[p3-1], bytes[p4-1]})
	return result
}

func GetFloat32(bytes []byte, p1 int, p2 int, p3 int, p4 int) (result float32) {
	result = byteToFloat32([]byte{bytes[p1-1], bytes[p2-1], bytes[p3-1], bytes[p4-1]})
	return result
}

func GetFloat64(bytes []byte, p1 int, p2 int, p3 int, p4 int, p5 int, p6 int, p7 int, p8 int) (result float64) {
	result = byteToFloat64([]byte{bytes[p1-1], bytes[p2-1], bytes[p3-1], bytes[p4-1], bytes[p5-1], bytes[p6-1],
		bytes[p7-1], bytes[p8-1]})
	log.Println(result)
	return result
}

func GetAscii(bytes []byte) string {
	return byteToAscii(bytes)
}

func GetHex(bytes []byte) []string {
	var resul []string
	for i := 1; i < len(bytes)-2; i += 2 {
		resul = append(resul, byteToHex([]byte{bytes[i-1], bytes[i]}))
	}
	return resul
}
