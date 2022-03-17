package main

import "rut955_openwrt/internal/modbus_rut"

func main() {
	//deviceDataChan := make(chan string, 1000)
	//modbus_rut.Start(deviceDataChan, "\"modbus_tcp_config.json\"")
	//modbus_rut.Test()
	arr := []byte{168,37,178,134,187,42,64,254,64,147,170,48,239,119,242,128}
	modbus_rut.GetFloat64(arr,1,2,3,4,5,6,7,8)
	modbus_rut.GetFloat64(arr,7,8,5,6,3,4,1,2)
}
