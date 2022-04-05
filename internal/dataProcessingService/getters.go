package dataProcessingService

import (
	"fmt"
	"github.com/fatih/set"
	"os/exec"
	"strconv"
	"strings"
	"time"
)


func getDeviceData(dataSourceChan chan string) [][]string {
	var dataList []string
	dataSet := set.New(set.ThreadSafe)
	for {
		select {
		case data := <-dataSourceChan:
			dataList = strings.Split(data,",")
			for _, s := range dataList {
				dataSet.Add(s)
			}

		default:
			dataList = nil
			for _, data := range dataSet.List() {
				dataList = append(dataList, fmt.Sprintf("%v",data))
			}
			return makeSlices(100, dataList)
		}
	}
}

func getDateTime() string {
	out, err := exec.Command("gpsctl", "-e").Output()
	if err != nil {
		out = []byte(time.Now().Format("2006-01-02 15:04:05"))
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
	out, err := exec.Command("gpsctl", "-u").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA"
	}
	
	strOut := strings.Replace(string(out), "\n", "", -1)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return fmt.Sprintf("%.3f", fltOut)
}

func getSats() string {
	out, err := exec.Command("gpsctl", "-p").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA"
	}

	strOut := strings.Replace(string(out), "\n", "", -1)
	return strOut
}

func getHeight() string {
	return "NA"
}

func getCourse() string {
	out, err := exec.Command("gpsctl", "-g").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA"
	}

	strOut := strings.Replace(string(out), "\n", "", -1)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return fmt.Sprintf("%.0f", fltOut)
}

func getSpeed() string {
	out, err := exec.Command("gpsctl", "-v").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA"
	}

	strOut := strings.Replace(string(out), "\n", "", -1)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return fmt.Sprintf("%.0f", fltOut)
}

func getLat() string {
	out, err := exec.Command("gpsctl", "-i").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}

	strOut := strings.Replace(string(out), "\n", "", -1)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	if fltOut == 0 {
		return "NA;NA"
	}
	if fltOut > 0 {
		return fmt.Sprintf("%.4f;N", fltOut*100)
	}
	return fmt.Sprintf("%.4f;S", fltOut*-100)

}

func getLon() string {
	out, err := exec.Command("gpsctl", "-x").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}
	strOut := strings.Replace(string(out), "\n", "", -1)
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	if fltOut == 0 {
		return "NA;NA"
	}
	if fltOut > 0 {
		return fmt.Sprintf("0%.4f;E", fltOut*100)
	}
	return fmt.Sprintf("0%.4f;W", fltOut*-100)
}
