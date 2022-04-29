package custom

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Custom struct {
	List []List `json:"list"`
}
type List struct {
	Mode string `json:"mode"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func Start(dataSourceChan chan string, path string) {
	defer log.Print("custom ds - done")
	log.Print("connect to custom data source")

	customConfig := getCustomConfig(path)

	for _, customModule := range customConfig.List {
		go customHandler(customModule, dataSourceChan)
	}
}

func customHandler(module List, dataSourceChan chan string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("something goes wrong in custom module.\nmodule - %v\nMsg > %v", module, r)
		}
	}()

	var address string
	switch module.Mode {
	case "auto":
		address = fmt.Sprintf("%v:%s", getOutboundIP(), module.Port)
	case "manual":
		address = fmt.Sprintf("%v:%s", module.IP, module.Port)
	}
	serverConnection, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicf("listenService error msg> %v",err)
	}
	defer serverConnection.Close()

	listenService(serverConnection, dataSourceChan)
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal("get ip error", err)
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Println("cant close connection")
		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func setDefaultCustomConfig() *Custom {
	return &Custom{
		List: []List{
			{
				Mode: "auto",
				IP:   "none",
				Port: "32211",
			},
		},
	}
}

func getCustomConfig(path string) (cfg *Custom) {
	cfg = setDefaultCustomConfig()
	configFile, err := os.Open(path)
	if err != nil {
		log.Printf("Using defaults. Bad config path : %v", path)
		return nil
	}
	defer configFile.Close()
	v := json.NewDecoder(configFile)
	_ = v.Decode(cfg)
	return
}

func getConfig(path string) *json.Decoder {
	configFile, err := os.Open(path)
	if err != nil {
		log.Printf("Using defaults. Bad config path : %v", path)
		return nil
	}
	defer configFile.Close()
	v := json.NewDecoder(configFile)
	return v
}
