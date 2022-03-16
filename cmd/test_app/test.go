package main

import (
	"rut955_openwrt/internal/modbus_rut"
)

func main() {
	modbusData := make(chan []string, 20)
	modbus_rut.Start(modbusData)
}
