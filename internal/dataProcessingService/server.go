package dataProcessingService

import (
	"fmt"
	"log"
	"rut_wialon_gateway/internal/modules/custom"
	"rut_wialon_gateway/internal/modules/mqtt"
	"strings"
	"time"
)

// type serverConfig struct {
//	MqttConfig   string
//	ModbusConfig string
//}

type Config struct {
	DataSourceChannelSize int     `json:"data_source_channel_size"`
	TickerDefTime         float64 `json:"ticker_def_time"`
	SpeedCoefficient      float64 `json:"speed_coefficient"`
	CourseDiffTrigger     int     `json:"course_diff_trigger"`
}

type Module struct {
	Name             string `json:"name"`
	ModuleConfigPath string `json:"module_config_path"`
}

type ModulesConfig struct {
	Modules []Module `json:"modules"`
}

func Start(dataChan chan string, config *Config, modulesConfig *ModulesConfig, dataSourceChan chan map[string][]string) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("recover in data processing service. Panic > %v", r)
			if strings.Contains(fmt.Sprintf("%v", r), "FATAL") {
				log.Panicf("FATAL error in data processing service. Use painc. Reason: %v", r)
			}
		}
	}()
	log.Print("Data Processing Service start")
	done := make(chan string)
	modulesConfig.connectDataSourceModules(dataSourceChan)
	go config.dataToWialonModule(dataChan, dataSourceChan, done)

	if d := <-done; true {
		log.Panicf("Restart Data Processing Service. Reason: %v", d)
	}
}

func (config *Config) dataToWialonModule(dataChan chan string, dataSourceChan chan map[string][]string, done chan string) {
	defer func() {
		if r := recover(); r != nil {
			done <- fmt.Sprintf("error in sendToDataChan, need restart. Reason: %v", r)
			return
		}
		done <- "dataToWialonModule for timer down"
	}()
	for {
		// log.Print("w8 and try to send")
		sendTimer(time.Now(), config.TickerDefTime, config.SpeedCoefficient, getCourseInt(), config.CourseDiffTrigger)
		sendToDataChan(dataChan, dataSourceChan)
	}
}

func sendTimer(startTime time.Time, rate float64, speedCoef float64, startCourse int, courseDiffTrigger int) {
	defFinishTime := startTime.Add(time.Duration(rate/(1.0+(getSpeedF64()*speedCoef*0.01))) * time.Second)
	for defFinishTime.After(time.Now()) && diff(startCourse, getCourseInt()) < courseDiffTrigger {
		time.Sleep(time.Duration(500/(1+getSpeedF64()*speedCoef*0.01)) * time.Millisecond)
	}
}

func diff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func sendToDataChan(dataChan chan string, dataSourceChan chan map[string][]string) {
	paramsList := getDeviceData(dataSourceChan)
	// paramList = [[params_block1],[params_block2]]
	var attr = []string{
		getDateTime(), getLat(), getLon(), getSpeed(), getCourse(), getHeight(), getSats(), getHdop(), getInputs(),
		getOutputs(), getAdc(), getIbutton(),
	}
	dataType := "D"
	if len(paramsList) == 0 {
		paramsList = []string{"NA"}
	}
	dataChan <- convertDataToSend(dataType, attr, paramsList)
	// log.Print("gata send to wialon module correct")
	// for
	// for _, params := range paramsList {
	// 	params = remove(params, "")
	// 	dataChan <- convertDataToSend(dataType, attr, params)
}

// }

// func remove(s []string, r string) []string {
// 	for i, v := range s {
// 		if v == r {
// 			return remove(append(s[:i], s[i+1:]...), r)
// 		}
// 	}
// 	return s
// }

func (config *ModulesConfig) connectDataSourceModules(dataSourceChan chan map[string][]string) {
	select {
	case d := <-dataSourceChan:
		dataSourceChan <- d
		return
	case <-time.After(16 * time.Second):
		for _, module := range config.Modules {
			log.Printf("start %v", module.Name)
			startModule(module, dataSourceChan)
		}
	}
}

func startModule(module Module, dataSourceChan chan map[string][]string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recover painc from module. Panic msg > %v", r)
		}
	}()
	switch module.Name {
	// case "modbus":
	// 	modbus_rut.Start(dataSourceChan, module.ModuleConfigPath)
	case "mqtt":
		mqtt.Start(dataSourceChan, module.ModuleConfigPath)
	case "custom":
		custom.Start(dataSourceChan, module.ModuleConfigPath)
	default:
		log.Printf("module %s not found", module.Name)
	}
}
