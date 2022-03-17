package main

import (
	"rut955_openwrt/internal/modbus_rut"
)

func main() {
	//deviceDataChan := make(chan string, 1000)
	//modbus_rut.Start(deviceDataChan, "\"modbus_tcp_config.json\"")
	modbus_rut.Test()
}
