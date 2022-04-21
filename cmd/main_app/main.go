package main

import (
	"encoding/json"
	"log"
	"os"
	"rut_wialon_gateway/internal/dataProcessingService"
	"rut_wialon_gateway/internal/wialonClient"
	"sync"
	"time"
)

type RutPathsConfig struct {
	WialonClientConfigPath          string `json:"wialon_client_config_path"`
	DataProcessingServiceConfigPath string `json:"data_processing_service_config_path"`
	ModulesConfigPath               string `json:"modules_config_path"`
	LogFilePath                     string `json:"log_file_path"`
	DataChannelSize                 int    `json:"data_channel_size"`
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("critical or unexpected panic > %v", r)
		}
		log.Print("main program ended. RWG_app_controller should restart.")
	}()

	setupLogger("/tmp/RWG_app_buffer/log.log")

	log.Print("start launch loop")
	configPath := "/overlay/rut_wialon_gateway/APP_PATHS.json" // path to main app config file
	for {
		launch(configPath)
		log.Print("restarting all service")
	}
}

func setupLogger(path string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("using defaults. panic > %v", r)
			logFile, err := os.OpenFile("/tmp/RWG_app_buffer/log.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return
			}
			log.SetOutput(logFile)
			log.SetFlags(log.Lshortfile | log.LstdFlags)
			log.Print("logger start")
		}
	}()

	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panicf("bad log path : %v", path)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Print("logger start")
}

func launch(path string) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	rutConfigPaths := getRutConfigPaths(path)
	wialonClientConfig, dataPSConfig, dataPSModulesConfig := getRutConfig(rutConfigPaths)

	dataChan := make(chan string, rutConfigPaths.DataChannelSize) // main data channel. Link dps-wialonClient

	log.Print("start main routines (threads)")
	// start wialon client thread
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("Critical error in wialon client. RWG_app_controller should restart. Error msg> %v", r)
			}
		}()
		defer wg.Done()
		for {
			startWialonClient(dataChan, wialonClientConfig)
			time.Sleep(10 * time.Second)
			log.Print("restarting wialon client process")
		}
	}()
	// start dps thread
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("Critical error in data processing service. RWG_app_controller should restart. Error msg> %v", r)
			}
		}()
		defer wg.Done()
		dataSourceChan := make(chan string, dataPSConfig.DataSourceChannelSize) // sub data channel. Link dataSources-converter
		for {
			startDataProcessingService(dataChan, dataPSConfig, dataPSModulesConfig, dataSourceChan)
			time.Sleep(10 * time.Second)
			log.Print("restarting data processing service process")
		}
	}()
	log.Print("wait for routines")
	wg.Wait()
}

func startDataProcessingService(
	dataChan chan string,
	dataPSConfig *dataProcessingService.Config,
	dataPSModulesConfig *dataProcessingService.ModulesConfig,
	dataSourceChan chan string,
) {
	dataProcessingService.Start(dataChan, dataPSConfig, dataPSModulesConfig, dataSourceChan)
}

func startWialonClient(dataChan chan string, wialonConfig *wialonClient.Config) {
	wialonClient.Start(dataChan, wialonConfig)
}

func getRutConfigPaths(path string) (cfg *RutPathsConfig) {
	cfg = setDefaultRutGatewayConfig()
	_ = getConfig(path)(&cfg)
	return
}

func getRutConfig(paths *RutPathsConfig) (
	wialonConfig *wialonClient.Config,
	dPSConfig *dataProcessingService.Config,
	dPSModulesConfig *dataProcessingService.ModulesConfig,
) {
	wialonConfig = getWialonConfig(paths.WialonClientConfigPath)
	dPSConfig = getDPSConfig(paths.DataProcessingServiceConfigPath)
	dPSModulesConfig = getDPSModulesConfig(paths.ModulesConfigPath)
	return
}

func getDPSModulesConfig(path string) (cfg *dataProcessingService.ModulesConfig) {
	cfg = setDefaultDPSModulesConfig()
	_ = getConfig(path)(&cfg)
	return
}

func getDPSConfig(path string) (cfg *dataProcessingService.Config) {
	cfg = setDefaultDataProcessingServiceConfig()
	_ = getConfig(path)(&cfg)
	return
}

func getWialonConfig(path string) (cfg *wialonClient.Config) {
	cfg = setDefaultWialonClientConfig()
	_ = getConfig(path)(&cfg)
	return
}

func getConfig(path string) func(v interface{}) error {
	configFile, err := os.Open(path)
	if err != nil {
		log.Printf("Using defaults. Bad config path : %v", path)
		return nil
	}
	defer configFile.Close()
	v := json.NewDecoder(configFile)
	return v.Decode
}

func setDefaultDPSModulesConfig() *dataProcessingService.ModulesConfig {
	return &dataProcessingService.ModulesConfig{
		Modules: []dataProcessingService.Module{{
			Name:             "mqtt",
			ModuleConfigPath: "/overlay/rut_wialon_gateway/MODULE_MQTT.json",
		}},
	}
}

func setDefaultDataProcessingServiceConfig() *dataProcessingService.Config {
	return &dataProcessingService.Config{
		DataSourceChannelSize: 1000,
		TickerTime:            10,
	}
}

func setDefaultWialonClientConfig() *wialonClient.Config {
	return &wialonClient.Config{
		WialonServerAddress: "10.0.0.2:11114",
		ConnectionType:      "tcp",
		DataBufferPath:      "/tmp/RWG_app_buffer/buffer.buf",
		Login:               "111222333444555",
		Password:            "",
	}
}

func setDefaultRutGatewayConfig() *RutPathsConfig {
	return &RutPathsConfig{
		WialonClientConfigPath:          "/overlay/rut_wialon_gateway/CFG_wialon_client.json",
		DataProcessingServiceConfigPath: "/overlay/rut_wialon_gateway/CFG_data_processing_service.json",
		ModulesConfigPath:               "/overlay/rut_wialon_gateway/MODULES_LIST.json",
		LogFilePath:                     "/tmp/RWG_app_buffer/log.log",
		DataChannelSize:                 50,
	}
}
