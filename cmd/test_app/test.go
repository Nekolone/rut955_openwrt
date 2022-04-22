package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//type Map struct {
//	m  map[string]interface{}
//	t  flatmap.Tokenizer
//	re *regexp.Regexp
//}

func main() {

	jsonStr := `{
        "fruits" : {
            "a": "apple",
            "b": "banana"
        },
        "colors" : {
            "r": "red",
            "g": "[1,2,3,4,5,6]"
        }
    }`

	data := make(map[string]interface{})
	res := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &data)
	Flatten2("", data, res)
	log.Print(res)
	fmt.Printf("New: %v\n", strings.ReplaceAll(fmt.Sprint(res["colors.g"]),"g", "n"))

	//fmt.Printf("Original: %v\n", Flatten(data))
}
func Flatten(m map[string]interface{}) map[string]interface{} {
	o := map[string]interface{}{}
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		case []interface{}:
			for i := 0; i < len(child); i++ {
				o[k+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			o[k] = v
		}
	}
	return o
}

func Flatten2(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	if len(prefix) > 0 {
		prefix += "."
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			Flatten2(prefix+k, child, dest)
		case []interface{}:
			for i := 0; i < len(child); i++ {
				dest[prefix+k+"."+strconv.Itoa(i)] = child[i]
			}
		default:
			dest[prefix+k] = v
		}
	}
}

//
//func anotherFunc() {
//	for {
//		log.Print("q")
//	}
//}
//
//func printInRout(ch chan int) {
//	defer func() {
//	    if r := recover(); r != nil {
//	        log.Panicf("WOOPS > %v", r)
//	    }
//	}()
//	i:=0
//	for {
//		log.Print("next")
//		ch <- i
//		i++
//		log.Panicf("hihihaha")
//	}
//}

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
