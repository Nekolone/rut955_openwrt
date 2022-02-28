package server

import (
	"log"
	"time"
)

func Start(networkStatus *string, dataChan chan string, iterChan chan string) {
	deviceDataChan := make(chan string)
	go getDataFromDevice(deviceDataChan)
	recTimer := time.NewTicker(time.Minute)
	log.Println("server start")
	for range recTimer.C {
		dataChan <- convertDeviceData(deviceDataChan)

	}

}

func convertDeviceData(deviceDataChan chan string) string {
	msg := "string"
	return msg
}

func getDataFromDevice(deviceDataChan chan string) {

}
