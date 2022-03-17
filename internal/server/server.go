package server

import (
	"log"
	"net"
	"rut955_openwrt/internal/modbus_rut"
	"time"
)

func Start(serverListenPort string, dataChan chan string, modbusConfigPath string) {
	log.Println("ListenServer start")

	deviceDataChan := make(chan string, 1000)

	go getDataFromDevices(serverListenPort, deviceDataChan, modbusConfigPath)

	recTimer := time.NewTicker(time.Second * 10)
	for range recTimer.C {
		sendToDataChan(dataChan, deviceDataChan)
	}
}

func sendToDataChan(dataChan chan string, deviceDataChan chan string) {
	dataList := getDeviceData(deviceDataChan)
	//data = ["","","","","","","","","","","",""]
	//msg = ["#type#params","#type#params","#type#params"]
	attr := []string{
		getDateTime(), getLat(), getLon(), getSpeed(), getCourse(), getHeight(), getSats(), getHdop(), getInputs(),
		getOutputs(), getAdc(), getIbutton(),
	}
	dataType := "D"
	for _, params := range dataList {
		dataChan <- convertDataToSend(dataType, attr, params)
	}
}

func getDataFromDevices(port string, deviceDataChan chan string, modbusConfigPath string) {
	modbus_rut.Start(deviceDataChan, modbusConfigPath)

	serverConnection, err := net.Listen("tcp", string(getOutboundIP())+port)
	if err != nil {
		log.Fatal("listenService error")
	}
	defer serverConnection.Close()

	listenService(serverConnection, deviceDataChan)
}
