package main

import (
	"log"
	"rut955_openwrt/internal/client"
	"rut955_openwrt/internal/server"
	"sync"
)

type config struct {
	connectionType string
	clientIp       string
}

func main() {
	configPath := "config_path"
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
	if err == nil {
		log.Println("config path error, using default config")
		Config = config{
			connectionType: "tcp",
			clientIp:       "192.168.57.161:12332",
		}
	}

	go func(wgp *sync.WaitGroup) {
		err := startClient(&Config, dataChan)
		defer wgp.Done()
		if err != nil {
			log.Fatal("client routine error")
		}
		log.Println("client stopped")
	}(&wg)

	go func(wgp *sync.WaitGroup) {
		err := startServer(dataChan)
		defer wgp.Done()
		if err != nil {
			log.Fatal("server routine error")
		}
		log.Println("server stopped")
	}(&wg)

	wg.Wait()
	return nil
}

func getConfig(path string) (config, error) {
	return config{}, nil
}

func startServer(dataChan chan string) error {
	server.Start(dataChan)
	return nil
}

func startClient(conf *config, dataChan chan string) error {
	networkStatus := "start"
	clientConnection, tcpAddr := client.Start(conf.clientIp, conf.connectionType, &networkStatus)
	go client.ReconnectingService(&tcpAddr, conf.connectionType, &clientConnection, &networkStatus)
	client.DataWorker(&networkStatus, &clientConnection, dataChan)
	//Сделать стоп и рестарт
	return nil
}

//func main() {
//	strEcho := "hi python server"
//	servAddr := "192.168.100.107:11113"
//
//	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
//	if err != nil {
//		log.Println("ResolveTCPAddr failed:", err.Error())
//		os.Exit(1)
//	}
//
//	conn, err := net.DialTCP("tcp", nil, tcpAddr)
//	if err != nil {
//		log.Println("Dial failed:", err.Error())
//		os.Exit(1)
//	}
//
//	_, err = conn.Write([]byte(strEcho))
//	if err != nil {
//		log.Println("Write to server failed:", err.Error())
//		os.Exit(1)
//	}
//
//	log.Println("write to server = ", strEcho)
//
//	reply := make([]byte, 1024)
//
//	_, err = conn.Read(reply)
//	if err != nil {
//		log.Println("Write to server failed:", err.Error())
//		os.Exit(1)
//	}
//
//	log.Println("reply from server=", string(reply))
//
//	err = conn.Close()
//	if err != nil {
//		log.Fatal("Close error")
//		return
//	}
//
//}
