package server

import (
	"fmt"
	"log"
	"math"
	"net"
	"os/exec"
	"strconv"
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
		getDateTime(), getLat(), getLon(), getSpeed(), getCourse(), getHeight(), getSats(), getHdop(), getInputs(),
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
	msg := ""
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

func getDateTime() string {
	out, err := exec.Command("gpsctl -e").Output()
	if err != nil {
		time.Now().Format("2006-01-02 15:04:05")
	}
	return string(out[8:10]) + string(out[5:7]) + string(out[2:4]) + ";" + string(out[11:13]) + string(out[14:16]) + string(out[17:19])
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
	out, err := exec.Command("gpsctl -u").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}
	strOut := string(out)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return fmt.Sprintf("%.3f", fltOut)
}

func getSats() string {
	out, err := exec.Command("gpsctl -p").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}
	strOut := string(out)
	return strOut
}

func getHeight() string {
	return "NA"
}

func getCourse() string {
	out, err := exec.Command("gpsctl -g").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}
	strOut := string(out)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	intOut := int(math.Round(fltOut))
	return fmt.Sprintf("%d", intOut)
}

func getSpeed() string {
	out, err := exec.Command("gpsctl -v").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}
	strOut := string(out)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	intOut := int(math.Round(fltOut))
	return fmt.Sprintf("%d", intOut)
}

func getLat() string {
	out, err := exec.Command("gpsctl -i").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}
	strOut := string(out)
	log.Println(strOut)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	if fltOut > 0 {
		return fmt.Sprintf("%f;N", fltOut*100)
	}
	return fmt.Sprintf("%f;S", fltOut*-100)

}

func getLon() string {
	out, err := exec.Command("gpsctl -x").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}
	strOut := string(out)
	log.Println(strOut)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	if fltOut > 0 {
		return fmt.Sprintf("0%f;E", fltOut*100)
	}
	return fmt.Sprintf("0%f;W", fltOut*-100)
}
