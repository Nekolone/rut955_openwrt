package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rut955_openwrt/internal/wialonClient"
	"rut955_openwrt/internal/dataProcessingService"
	"sync"
)

type RutPathsConfig struct {
	WialonClientConfigPath          string `json:"wialon_client_config_path"`
	DataProcessingServiceConfigPath string `json:"data_processing_service_config_path"`
	ModulesConfigPath               string `json:"modules_config_path"`
}

func main() {
	configPath := "rut_gateway_config_paths.json"
	err := launch(configPath)
	if err != nil {
		log.Fatal("launch error")
		return
	}
	log.Println("service finished")
}

func launch(path string) (err error) {
	log.Println("enter launch")

	dataChan := make(chan string, 50)

	wg := sync.WaitGroup{}
	wg.Add(2)

	log.Println("launch - get cfgs")

	var rutConfigPaths *RutPathsConfig
	rutConfigPaths, err = getRutConfigPaths(path)
	if err != nil {
		log.Println("RutPathsConfig path error, using default RutPathsConfig")
		rutConfigPaths = setDefaultRutGatewayConfig()
	}
	wialonClientConfig, dataPSConfig, dataPSModulesConfig := getRutConfig(rutConfigPaths)

	log.Println("launch - start routines")

	go func(wgp *sync.WaitGroup) {
		defer wgp.Done()
		if err = startWialonClient(dataChan, wialonClientConfig); err != nil {
			log.Fatal("wialon client routine error")
		}
		log.Println("wialon client stopped")
	}(&wg)

	go func(wgp *sync.WaitGroup) {
		defer wgp.Done()
		if err = startDataProcessingService(dataChan, dataPSConfig, dataPSModulesConfig); err != nil {
			log.Fatal("data processing routine error")
		}
		log.Println("dps stopped")
	}(&wg)

	wg.Wait()
	log.Println("launch - routines end")
	return nil
}

func startDataProcessingService(dataChan chan string, dataPSConfig *dataProcessingService.Config,
	dataPSModulesConfig *dataProcessingService.ModulesConfig) error {

	dataProcessingService.Start(dataChan, dataPSConfig, dataPSModulesConfig)
	return nil
}

func startWialonClient(dataChan chan string, wialonConfig *wialonClient.Config) error {

	wialonClient.Start(dataChan, wialonConfig)
	return nil
}

func getRutConfig(paths *RutPathsConfig) (*wialonClient.Config, *dataProcessingService.Config, *dataProcessingService.ModulesConfig) {
	wialonConfig, err := getWialonConfig(paths.WialonClientConfigPath)
	if err != nil {
		wialonConfig = setDefaultWialonClientConfig()
	}
	dPSConfig, err := getDPSConfig(paths.DataProcessingServiceConfigPath)
	if err != nil {
		dPSConfig = setDefaultDataProcessingServiceConfig()
	}
	dPSModulesConfig, err := getDPSModulesConfig(paths.ModulesConfigPath)
	if err != nil {
		dPSModulesConfig = setDefaultDPSModulesConfig()
	}
	return wialonConfig, dPSConfig, dPSModulesConfig
}

func getDPSModulesConfig(path string) (*dataProcessingService.ModulesConfig, error) {
	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return &dataProcessingService.ModulesConfig{}, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	var cfg dataProcessingService.ModulesConfig
	_ = jsonParser.Decode(&cfg)
	return &cfg, nil
}

func getDPSConfig(path string) (*dataProcessingService.Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return &dataProcessingService.Config{}, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	var cfg dataProcessingService.Config
	_ = jsonParser.Decode(&cfg)
	return &cfg, nil
}

func getWialonConfig(path string) (*wialonClient.Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return &wialonClient.Config{}, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	var cfg wialonClient.Config
	_ = jsonParser.Decode(&cfg)
	return &cfg, nil
}

func setDefaultDPSModulesConfig() *dataProcessingService.ModulesConfig {
	return &dataProcessingService.ModulesConfig{
		Modules: []dataProcessingService.Module{{
			Name: "mqtt",
			ModuleConfigPath: "module_mqtt_config.json",
		}},
	}
}

func setDefaultDataProcessingServiceConfig() *dataProcessingService.Config {
	return &dataProcessingService.Config{
		DataSourceChannelSize: 1000,
		TickerTime: 10,
	}
}

func setDefaultWialonClientConfig() *wialonClient.Config {
	return &wialonClient.Config{

	}
}

func setDefaultRutGatewayConfig() *RutPathsConfig {
	return &RutPathsConfig{
		WialonClientConfigPath:          "rut_wialon_client_config.json",
		DataProcessingServiceConfigPath: "rut_data_processing_service_config.json",
		ModulesConfigPath:               "rut_modules_config.json",
	}
}

func getRutConfigPaths(path string) (*RutPathsConfig, error) {

	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return &RutPathsConfig{}, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	var cfg RutPathsConfig
	_ = jsonParser.Decode(&cfg)
	return &cfg, nil
}
