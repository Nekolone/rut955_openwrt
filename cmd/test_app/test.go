package main

import (
	"github.com/fatih/set"
	"log"
)

func main() {
	//deviceDataChan := make(chan string, 1000)
	//modbus_rut.Start(deviceDataChan, "\"modbus_tcp_config.json\"")
	//modbus_rut.Test()
	//	arr := []byte{168,37,178,134,187,42,64,254,64,147,170,48,239,119,242,128}
	//	modbus_rut.GetFloat64(arr,1,2,3,4,5,6,7,8)
	//	modbus_rut.GetFloat64(arr,7,8,5,6,3,4,1,2)
	//}

	s := set.New(set.ThreadSafe) // thread safe version
	s.Add("istanbul")
	s.Add("istanbul1")
	s.Add("istanbul13")
	s.Add("istanbul134")
	log.Println(s)
	a := s.List()
	log.Println(a)
	for _, i := range a {
		log.Println(i)
	}
}
