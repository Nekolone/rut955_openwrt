package modbus_rut

func GetSigInt8(bytes []byte) []int8 {
	var resul []int8
	for _, d := range bytes {
		resul = append(resul, int8(d))
	}
	return resul
}

func GetUnsInt8(bytes []byte) []uint8 {
	return bytes
}

func GetSigInt16(bytes []byte, p1 int, p2 int) []int16 {
	var resul []int16
	for i := 1; i <= len(bytes)-2; i += 2 {
		resul = append(resul, byteToInt16([]byte{bytes[i-(2-p1)], bytes[i-(2-p2)]}))
	}
	return resul
}

func GetUnsInt16(bytes []byte, p1 int, p2 int) []uint16 {
	var resul []uint16
	for i := 1; i <= len(bytes)-2; i += 2 {
		resul = append(resul, byteToUint16([]byte{bytes[i-(2-p1)], bytes[i-(2-p2)]}))
	}
	return resul
}

func GetSigInt32(bytes []byte, p1 int, p2 int, p3 int, p4 int) []int32 {
	var resul []int32
	for i := 1; i <= len(bytes)-2; i += 2 {
		resul = append(resul, byteToInt32([]byte{bytes[i-(4-p1)], bytes[i-(4-p2)], bytes[i-(4-p3)], bytes[i-(4-p4)]}))
	}
	return resul
}

func GetUnsInt32(bytes []byte, p1 int, p2 int, p3 int, p4 int) []uint32 {
	var resul []uint32
	for i := 1; i <= len(bytes)-2; i += 2 {
		resul = append(resul, byteToUint32([]byte{bytes[i-(4-p1)], bytes[i-(4-p2)], bytes[i-(4-p3)], bytes[i-(4-p4)]}))
	}
	return resul
}

func GetFloat32(bytes []byte, p1 int, p2 int, p3 int, p4 int) []float32 {
	var resul []float32
	for i := 3; i <= len(bytes)-4; i += 4 {
		resul = append(resul, byteToFloat32([]byte{bytes[i-(4-p1)], bytes[i-(4-p2)], bytes[i-(4-p3)], bytes[i-(4-p4)]}))
	}
	return resul
}

func GetFloat64(bytes []byte, p1 int, p2 int, p3 int, p4 int, p5 int, p6 int, p7 int, p8 int) []float64 {
	var resul []float64
	for i := 7; i <= len(bytes)-8; i += 8 {
		resul = append(resul, byteToFloat64([]byte{bytes[i-(8-p1)], bytes[i-(8-p2)], bytes[i-(8-p3)], bytes[i-(8-p4)],
			bytes[i-(8-p5)], bytes[i-(8-p6)], bytes[i-(8-p7)], bytes[i-(8-p8)]}))
	}

	return resul
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
