package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rut955_openwrt/internal/client"
	"rut955_openwrt/internal/server"
	"sync"
)

type config struct {
	DeviceId         string `json:"device_id"`
	DevicePass       string `json:"device_pass"`
	ConnectionType   string `json:"connection_type"`
	ClientIp         string `json:"client_ip"`
	BufferPath       string `json:"buffer_path"`
	ServerListenPort string `json:"server_listen_port"`
	ModbusConfigPath string `json:"modbusConfigPath"`
}

func main() {
	configPath := "rut_config.json"
	err := launch(configPath)
	if err != nil {
		log.Fatal("launch error")
		return
	}
	log.Println("service finished")
}

func launch(path string) error {

	dataChan := make(chan string, 20)

	wg := sync.WaitGroup{}
	wg.Add(2)

	Config, err := getConfig(path)
	if err != nil {
		log.Println("config path error, using default config")
		Config = config{
			DeviceId:         "00000000001",
			DevicePass:       "passwd",
			ConnectionType:   "tcp",
			ClientIp:         "192.168.100.107:11113",
			BufferPath:       "buffer.buf",
			ServerListenPort: ":11111",
			ModbusConfigPath: "modbus_tcp_config.json",
		}
	}

	go func(wgp *sync.WaitGroup) {
		defer wgp.Done()
		err = startClient(&Config, dataChan)
		if err != nil {
			log.Fatal("client routine error")
		}
		log.Println("client stopped")
	}(&wg)

	go func(wgp *sync.WaitGroup) {
		defer wgp.Done()
		err = startServer(&Config, dataChan)
		if err != nil {
			log.Fatal("server routine error")
		}
		log.Println("server stopped")
	}(&wg)

	wg.Wait()
	return nil
}

func getConfig(path string) (config, error) {
	var cfg config
	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return config{}, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	_ = jsonParser.Decode(&cfg)
	return cfg, nil
}

func startServer(conf *config, dataChan chan string) error {
	server.Start(conf.ServerListenPort, dataChan, conf.ModbusConfigPath)
	return nil
}

func startClient(conf *config, dataChan chan string) error {
	networkStatus := "start"

	clientConnection, tcpAddr := client.ConnectToServer(conf.ClientIp, conf.ConnectionType, &networkStatus,
		conf.DeviceId, conf.DevicePass)
	go client.ReconnectingService(&tcpAddr, conf.ConnectionType, &clientConnection, &networkStatus, conf.DeviceId,
		conf.DevicePass)
	client.DataWorker(&networkStatus, &clientConnection, dataChan, conf.BufferPath)
	//Сделать стоп и рестарт
	return nil
}
