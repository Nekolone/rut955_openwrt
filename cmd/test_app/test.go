package main

import (
	"fmt"
	"time"
)

func main() {
	for{
		fmt.Print("13")
		time.Sleep(2 * time.Second)
	}
}

//type RutPathsConfig struct {
//	WialonClientConfigPath          string `json:"wialon_client_config_path"`
//	DataProcessingServiceConfigPath string `json:"data_processing_service_config_path"`
//	ModulesConfigPath               string `json:"modules_config_path"`
//	LogFilePath                     string `json:"log_file_path"`
//	DataChannelSize                 int    `json:"data_channel_size"`
//}

//func main() {
//	defer func() {
//		if r := recover(); r != nil {
//			log.Printf("Recovered in f > %v", r)
//		}
//	}()
//
//	logFile, err := os.OpenFile("log.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
//	if err != nil {
//		logFile, err = os.OpenFile("log.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
//	}
//	log.SetOutput(logFile)
//	log.SetFlags(log.Lshortfile | log.LstdFlags)
//
//	log.Print("logger start")
//	rutConfigPaths := getRutConfigPaths("123")
//	log.Print(*rutConfigPaths)
//	log.Print(rutConfigPaths.DataChannelSize)
//	log.Print(rutConfigPaths.ModulesConfigPath)
//	log.Print(rutConfigPaths.DataProcessingServiceConfigPath)
//	log.Print(rutConfigPaths.WialonClientConfigPath)
//
//}



//func setDefaultRutGatewayConfig() *RutPathsConfig {
//	return &RutPathsConfig{
//		WialonClientConfigPath:          "/overlay/rut_wialon_gateway/CFG_wialon_client.json",
//		DataProcessingServiceConfigPath: "/overlay/rut_wialon_gateway/CFG_data_processing_service.json",
//		ModulesConfigPath:               "/overlay/rut_wialon_gateway/MODULES_LIST.json",
//		LogFilePath:                     "/tmp/RWG_app_buffer/log.log",
//	}
//}
//
//func getRutConfigPaths(path string) (cfg *RutPathsConfig) {
//	defer func() {
//		if r := recover(); r != nil {
//			log.Printf("Recover {%v}. Using defaults", r)
//			cfg = setDefaultRutGatewayConfig()
//		}
//	}()
//
//
//	configFile, err := os.Open(path)
//	if err != nil {
//		log.Panicf("bad rut cfg path : %v", err.Error())
//	}
//	defer configFile.Close()
//
//	jsonParser := json.NewDecoder(configFile)
//	_ = jsonParser.Decode(&cfg)
//	return cfg
//}


//
//func we() {
//	// LOG_FILE := "/overlay/wialon_rut955_gateway/test.log"
//	LOG_FILE := "test.log"
//
//	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
//	if err != nil {
//		log.Printf("cant open log file. ERR > %v", err)
//	}
//
//	defer logFile.Close()
//
//	log.SetOutput(logFile)
//
//	log.SetFlags(log.Lshortfile | log.LstdFlags)
//
//	log.Println("test")
//
//	var wg sync.WaitGroup
//
//	wg.Add(1)
//
//	in_func_test()
//	go func() {
//		defer wg.Done()
//		log.Println("go")
//	}()
//
//	time.Sleep(1)
//	wg.Wait()
//	time.Sleep(1)
//
//}
//
//func in_routine_test(wg sync.WaitGroup) {
//	log.Println("in routine test")
//	wg.Done()
//}
//func in_func_test() {
//	log.Println("in func test")
//}
