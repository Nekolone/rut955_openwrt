package dataProcessingService

import (
	"fmt"
	"log"
	"rut_wialon_gateway/internal/modules/custom"
	"rut_wialon_gateway/internal/modules/modbus_rut"
	"rut_wialon_gateway/internal/modules/mqtt"
	"strings"
	"time"
)

// type serverConfig struct {
//	MqttConfig   string
//	ModbusConfig string
//}

type Config struct {
	DataSourceChannelSize int `json:"data_source_channel_size"`
	TickerTime            int `json:"ticker_time"`
}

type Module struct {
	Name             string `json:"name"`
	ModuleConfigPath string `json:"module_config_path"`
}

type ModulesConfig struct {
	Modules []Module `json:"modules"`
}

func Start(dataChan chan string, config *Config, modulesConfig *ModulesConfig, dataSourceChan chan string) {
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
		log.Panicf("Something gose wrong: %v", d)
	}
}

func (config *Config) dataToWialonModule(dataChan, dataSourceChan, done chan string) {
	defer func() {
		if r := recover(); r != nil {
			done <- fmt.Sprintf("error in sendToDataChan, need restart. Reason: %v", r)
			return
		}
		done <- "dataToWialonModule for timer down"
	}()
	for range time.NewTicker(time.Second * time.Duration(config.TickerTime)).C {
		sendToDataChan(dataChan, dataSourceChan)
	}
}

func sendToDataChan(dataChan, dataSourceChan chan string) {
	paramsList := getDeviceData(dataSourceChan)
	var attr = []string{
		getDateTime(), getLat(), getLon(), getSpeed(), getCourse(), getHeight(), getSats(), getHdop(), getInputs(),
		getOutputs(), getAdc(), getIbutton(),
	}
	dataType := "D"
	for _, params := range paramsList {
		params = remove(params, "")
		if len(params) == 0 {
			params = []string{"NA"}
		}
		dataChan <- convertDataToSend(dataType, attr, params)
	}
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return remove(append(s[:i], s[i+1:]...), r)
		}
	}
	return s
}

func (config *ModulesConfig) connectDataSourceModules(dataSourceChan chan string) {
	select {
	case d := <-dataSourceChan:
		dataSourceChan <- d
		return
	case <-time.After(16 * time.Second):
		for _, module := range config.Modules {
			startModule(module, dataSourceChan)
		}
	}
}

func startModule(module Module, dataSourceChan chan string) {
	defer recoverDPSDefaultPanic()

	switch module.Name {
	case "modbus":
		modbus_rut.Start(dataSourceChan, module.ModuleConfigPath)
	case "mqtt":
		mqtt.Start(dataSourceChan, module.ModuleConfigPath)
	case "custom":
		custom.Start(dataSourceChan, module.ModuleConfigPath)

	default:
		log.Printf("module %s not found", module.Name)
	}
}

func recoverDPSDefaultPanic() {
	if r := recover(); r != nil {
		log.Printf("Recover painc from module. Panic msg > %v", r)
	}
}
