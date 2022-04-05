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
	customConfig, err := getMqttConfig(path)
	if err != nil {
		log.Println("mqtt config path error, default config for local mqtt broker")
		customConfig = setDefaultCustomConfig()
	}

	for _, customModule := range customConfig.List {
		go customHandler(customModule, dataSourceChan)
	}
}

func customHandler(module List, dataSourceChan chan string) {
	var address string
	switch module.Mode {
	case "auto":
		address = fmt.Sprintf("%v:%s",getOutboundIP(),module.Port)
	case "manual":
		address = fmt.Sprintf("%v:%s",module.IP,module.Port)
	}
	serverConnection, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("listenService error")
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
				Mode:   "auto",
				IP: "none",
				Port: "32211",
			},
		},
	}
}

func getMqttConfig(path string) (*Custom, error) {
	var cfg Custom
	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return &Custom{}, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	_ = jsonParser.Decode(&cfg)
	return &cfg, nil
}
