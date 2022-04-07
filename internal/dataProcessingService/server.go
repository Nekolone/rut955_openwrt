package dataProcessingService

import (
	"log"
	"rut955_openwrt/internal/modules/custom"
	"rut955_openwrt/internal/modules/modbus_rut"
	"rut955_openwrt/internal/modules/mqtt"
	"time"
)

//type serverConfig struct {
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

func Start(dataChan chan string, config *Config, modulesConfig *ModulesConfig) {
	log.Println("Data Processing Service start")

	dataSourceChan := make(chan string, config.DataSourceChannelSize)

	modulesConfig.connectDataSourceModules(dataSourceChan)

	for range time.NewTicker(time.Second * time.Duration(config.TickerTime)).C {
		sendToDataChan(dataChan, dataSourceChan)
	}
}

func sendToDataChan(dataChan chan string, dataSourceChan chan string) {
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
	for _, module := range config.Modules {
		startModule(&module, dataSourceChan)
	}
}

func startModule(module *Module, dataSourceChan chan string) {
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
