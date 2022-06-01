package dataProcessingService

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getDeviceData(dataSourceChan chan map[string][]string) (res []string) {
	// var dataList []string
	// dataSourceChan = 'mqtt_mod1':{"time", "data"}
	//
	//res = map{
	//	"mqtt_mod1":{
	//		{time, data},
	//		{time, data},
	//		{time, data}
	//  },
	//   "mqtt_mod2":{
	//		{time, data},
	//		{time, data},
	//		{time, data}
	//  }
	//}
	// dataSet := set.New(set.ThreadSafe)
	dataList := make(map[string][][]string)
	for {
		select {
		case data := <-dataSourceChan:
			for k := range data {
				dataList[k] = append(dataList[k], data[k])
			}
			// dataSet.Add(data)
		default:
			for k := range dataList {
				// param := fmt.Sprintf("%v:string:%v", k, strings.Join(getParamData(dataList[k]),"%"))
				res = append(res, fmt.Sprintf("%v:list:%v", k, strings.Join(getParamData(dataList[k]), "%")))
			}
			// dataList = nil
			// for _, data := range dataSet.List() {
			// 	dataList = append(dataList, fmt.Sprintf("%v", data))
			// }
			// return makeSlices(100, dataList)
		}
	}
}

func getParamData(dataList [][]string) (res []string) {
	for _, v := range dataList {
		res = append(res, strings.Join(v, "*"))
	}
	return
}

// func getSubParamData(data []string) []string{

// }

func getDateTime() string {
	out, err := exec.Command("gpsctl", "-e").Output()
	if err != nil || bytes.Equal(out, []byte("1970-01-01 02:00:00")) {
		out = []byte(time.Now().Format("2006-01-02 15:04:05"))
	}
	log.Print("gata get correct")
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

	strOut := strings.ReplaceAll(string(out), "\n", "")
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return fmt.Sprintf("%.3f", fltOut)
}

func getSats() string {
	out, err := exec.Command("gpsctl", "-p").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA"
	}

	strOut := strings.ReplaceAll(string(out), "\n", "")
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

	strOut := strings.ReplaceAll(string(out), "\n", "")
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return fmt.Sprintf("%.0f", fltOut)
}

func getCourseInt() int {
	out, err := exec.Command("gpsctl", "-g").Output()
	if (err != nil) || (len(out) == 0) {
		return 0
	}

	strOut := strings.ReplaceAll(string(out), "\n", "")
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return int(fltOut)
}

func getSpeed() string {
	out, err := exec.Command("gpsctl", "-v").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA"
	}

	strOut := strings.ReplaceAll(string(out), "\n", "")
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return fmt.Sprintf("%.0f", fltOut)
}

func getSpeedF64() float64 {
	out, err := exec.Command("gpsctl", "-v").Output()
	if (err != nil) || (len(out) == 0) {
		return 0
	}

	strOut := strings.ReplaceAll(string(out), "\n", "")
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	return fltOut
}

func getLat() string {
	out, err := exec.Command("gpsctl", "-i").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}

	strOut := strings.ReplaceAll(string(out), "\n", "")
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
	strOut := strings.ReplaceAll(string(out), "\n", "")
	fltOut, _ := strconv.ParseFloat(strOut, 64)
	if fltOut == 0 {
		return "NA;NA"
	}
	if fltOut > 0 {
		return fmt.Sprintf("0%.4f;E", fltOut*100)
	}
	return fmt.Sprintf("0%.4f;W", fltOut*-100)
}
