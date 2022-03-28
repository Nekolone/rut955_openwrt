package dataProcessingService

import (
	"encoding/json"
	"log"
	"os"
	"rut955_openwrt/internal/modules/modbus_rut"
	"rut955_openwrt/internal/modules/mqtt"
	"time"
)

type serverConfig struct {
	MqttConfig   string
	ModbusConfig string
}

type Config struct {


}

type ModulesConfig struct {

}

func Start(dataChan chan string, config *Config, modulesConfig *ModulesConfig) {
	log.Println("ListenServer start")

	//config, err := getServerConfig("data_server_config.json")
	//if err != nil {
	//	config = serverConfig{
	//		MqttConfig:   "mqtt_config.json",
	//		ModbusConfig: "modbus_tcp_config.json",
	//	}
	//}

	dataSourceChan := make(chan string, 1000)

	modulesConfig.connectDataSourceModules(dataSourceChan)

	for range time.NewTicker(time.Second * 10).C {
		sendToDataChan(dataChan, dataSourceChan)
	}
}

func getServerConfig(path string) (config serverConfig, err error) {
	var configFile *os.File
	configFile, err = os.Open(path)
	if err != nil {
		log.Fatalf("cant open %s, err %v\n", path, err.Error())
		return
	}
	defer configFile.Close()
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatalf("cant decode %s, err %v\n", path, err.Error())
		return
	}
	return
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

func (config *ModulesConfig) connectDataSourceModules(deviceDataChan chan string) {

	modbus_rut.Start(deviceDataChan, config.ModbusConfig)
	mqtt.Start(deviceDataChan, config.MqttConfig)

	//serverConnection, err := net.Listen("tcp", string(getOutboundIP())+port)
	//if err != nil {
	//	log.Fatal("listenService error")
	//}
	//defer serverConnection.Close()
	//
	//listenService(serverConnection, deviceDataChan)
}
