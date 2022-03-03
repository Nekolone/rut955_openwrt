package server

import (
	"log"
	"net"
	"strings"
	"time"
)

func Start(dataChan chan string) {
	deviceDataChan := make(chan string, 1000)
	go getDataFromDevices(deviceDataChan)
	recTimer := time.NewTicker(time.Minute)
	log.Println("server start")
	for range recTimer.C {
		sendToDataChan(dataChan, deviceDataChan)

	}

}

func sendToDataChan(dataChan chan string, deviceDataChan chan string) {

	dataList := getDeviceData(deviceDataChan)
	//data = ["","","","","","","","","","","",""]
	//msg = ["#type#params","#type#params","#type#params"]
	attr := []string{
		getDate(), getTime(), getGps(), getSpeed(), getCourse(), getHeight(), getSats(), getHdop(), getInputs(),
		getOutputs(), getAdc(), getIbutton(),
	}
	for _, params := range dataList {
		log.Println(listToSrt(params, ","))
		dataChan <- "#D#" + strings.Join(attr, ";") + ";" + listToSrt(params, ",") + "\r\n"
	}

}

func listToSrt(params []string, delim string) string {
	if len(params) == 0 {
		return ""
	}
	var msg string
	for i := 0; i < len(params)-1; i++ {
		msg = params[i] + delim
	}
	return msg + params[len(params)-1]
}
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal("get ip error", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("cant close connection")
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func convertDeviceData(deviceDataChan chan string) []string {
	dataList := getDeviceData(deviceDataChan)
	//data = ["","","","","","","","","","","",""]
	//msg = ["#type#params","#type#params","#type#params"]
	attr := []string{
		getDate(), getTime(), getGps(), getSpeed(), getCourse(), getHeight(), getSats(), getHdop(), getInputs(),
		getOutputs(), getAdc(), getIbutton(),
	}
	var msg []string
	for _, ms := range dataList {
		msg = append(msg, "#D#"+strings.Join(attr, ";")+strings.Join(ms, ","))
	}

	return msg
}

func getDeviceData(deviceDataChan chan string) [][]string {
	var dataList []string
	for {
		select {
		case data := <-deviceDataChan:
			dataList = append(dataList, data)
		default:
			return makeSlices(100, dataList)
		}
	}

}

func makeSlices(i int, list []string) [][]string {
	if len(list) <= i {
		return [][]string{list}
	}
	return append(makeSlices(i, list[i:]), list[:i])
}

func getDataFromDevices(deviceDataChan chan string) {
	serverConnection, err := net.Listen("tcp", string(getOutboundIP())+":"+"16969")
	if err != nil {
		log.Fatal("listenService error")
	}
	defer serverConnection.Close()

	listenService(serverConnection, deviceDataChan)

}

func listenService(serverConnection net.Listener, deviceDataChan chan string) {
	for {
		deviceConnection, err := serverConnection.Accept()
		if err != nil {
			log.Println("listen service error")
		}
		go handleRequest(deviceConnection, deviceDataChan)
	}
}

func handleRequest(connection net.Conn, deviceDataChan chan string) {
	for {
		buffer := make([]byte, 1024)
		reqLen, err := connection.Read(buffer)
		if err != nil {
			log.Println("handle request error", err)
			return
		}
		if reqLen == 0 {
			continue
		}

		deviceDataChan <- string(buffer)
	}
}

func getIbutton() string {
	return "NA"
}

func getAdc() string {
	return "NA"
}

func getOutputs() string {
	return "NA"
}

func getInputs() string {
	return "NA"
}

func getHdop() string {
	return "NA"
}

func getSats() string {
	return "NA"
}

func getHeight() string {
	return "NA"
}

func getCourse() string {
	return "NA"
}

func getSpeed() string {
	return "NA"
}

func getGps() string {
	return "NA"
}

func getTime() string {
	return "NA"
}

func getDate() string {
	return "NA"
}
