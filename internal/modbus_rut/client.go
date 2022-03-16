package modbus_rut

import (
	"github.com/goburrow/modbus"
	"log"
	"os"
	"time"
)

type ModbusDevices struct {
	deviceList []Device
}

type Device struct {
	deviceIp string
	timeout  time.Duration
	requests []Request
}
type Request struct {
	slaveId      byte
	function     int
	startAddress uint16
	quantity     uint16
	registerMap  RegMap
}

type RegMap struct {
}

func Start(modbusData chan []string) {
	var devices ModbusDevices
	for _, device := range devices.deviceList {
		go deviceWorker(device, modbusData)
	}
}

func deviceWorker(device Device, modbusData chan []string) {
	handler := modbus.NewTCPClientHandler(device.deviceIp)
	handler.Timeout = device.timeout
	for {
		err := handler.Connect()
		if err != nil {
			log.Println(err)
			continue
		}
		RequestsService(device.requests, modbusData, handler)
	}
}

func RequestsService(requests []Request, modbusData chan []string, handler *modbus.TCPClientHandler) {
	defer handler.Close()
	var client modbus.Client
	var results []byte
	var err error
	for {
		for _, request := range requests {
			handler.SlaveId = request.slaveId
			client = modbus.NewClient(handler)
			switch request.function {
			case 1:
				continue
			case 2:
				continue
			case 3:
				results, err = client.ReadHoldingRegisters(request.startAddress, request.quantity)
				if err != nil {
					log.Printf("RequestsService err %v\n", err)
				}
			default:
				continue
			}
			modbusData <- modbusConvertService(results, request.registerMap)
		}
	}
}

func GetModbus321321ata() []byte {
	handler := modbus.NewTCPClientHandler("172.16.1.12:502")
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0x0C
	handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	if err != nil {
		log.Println(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadHoldingRegisters(0, 20)
	return results
}
