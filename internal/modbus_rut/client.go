package modbus_rut

import (
	"encoding/json"
	"github.com/goburrow/modbus"
	"log"
	"os"
	"time"
)

type ModbusDevices []struct {
	Name     string    `json:"name"`
	DeviceIp string    `json:"deviceIp"`
	Timeout  int64       `json:"timeout"`
	Requests []Request `json:"requests"`
}

type ModbusDevice struct {
	Name     string
	DeviceIp string
	Timeout  int64
	Requests []Request
}

type Request struct {
	SlaveId         byte   `json:"slaveId"`
	Function        int    `json:"function"`
	StartAddress    uint16 `json:"startAddress"`
	Quantity        uint16 `json:"quantity"`
	RegisterMap     []param `json:"params"`
}

type param struct {
	ParamId  string `json:"paramId"`
	ByteSize uint16 `json:"byteSize"`
	DataType string `json:"dataType"`
	ByteOrder []int `json:"byteOrder"`
}

func Test() {

	deviceList := getModbusRequestsConfig("modbus_tcp_config.json")



	for _, device := range deviceList {
		log.Println(device.Name)
		log.Println(device.DeviceIp)
		log.Println(device.Timeout)
		for _, request := range device.Requests {
			log.Println(request.SlaveId)
			log.Println(request.Function)
			log.Println(request.StartAddress)
			log.Println(request.Quantity)
			log.Println(request.RegisterMap)
			for _, p := range request.RegisterMap {
				log.Println("paramId", p.ParamId)
				log.Println("byteSize", p.ByteSize)
				log.Println("dataType", p.DataType)
				log.Println("byteOrder", p.ByteOrder)
			}
		}
	}
}


func Start(deviceDataChan chan string, path string) {

	deviceList := getModbusRequestsConfig(path)
	for _, device := range deviceList {
		go deviceWorker(ModbusDevice(device), deviceDataChan)
	}

}

func getModbusRequestsConfig(path string) (deviceList ModbusDevices) {
	var configFile *os.File
	var err error
	configFile, err = os.Open(path)
	if err != nil {
		log.Fatalf("cant open %s, err %v\n", path, err.Error())
		return nil
	}
	defer configFile.Close()
	err = json.NewDecoder(configFile).Decode(&deviceList)
	if err != nil {
		log.Fatalf("cant decode %s, err %v\n", path, err.Error())
		return nil
	}
	return
}

func deviceWorker(device ModbusDevice, deviceDataChan chan string) {
	handler := modbus.NewTCPClientHandler(device.DeviceIp)
	handler.Timeout = time.Duration(device.Timeout) * time.Second
	for {
		err := handler.Connect()
		if err != nil {
			log.Println(err)
			continue
		}
		RequestsService(device.Requests, deviceDataChan, handler)
	}
}

func RequestsService(requests []Request, deviceDataChan chan string, handler *modbus.TCPClientHandler) {
	defer handler.Close()
	var client modbus.Client
	var results []byte
	recTimer := time.NewTicker(time.Second*10)
	for range recTimer.C {
		for _, request := range requests {
			handler.SlaveId = request.SlaveId
			client = modbus.NewClient(handler)
			switch request.Function {
			case 1:
				results = errHandlerBA(client.ReadCoils(request.StartAddress, request.Quantity))
			case 2:
				results = errHandlerBA(client.ReadDiscreteInputs(request.StartAddress, request.Quantity))
			case 3:
				results = errHandlerBA(client.ReadHoldingRegisters(request.StartAddress, request.Quantity))
			case 4:
				results = errHandlerBA(client.ReadInputRegisters(request.StartAddress, request.Quantity))
			case 5:
				results = errHandlerBA(client.WriteSingleCoil(request.StartAddress, uint16(69)))
			case 6:
				results = errHandlerBA(client.WriteSingleRegister(request.StartAddress, uint16(69)))
			case 15:
				results = errHandlerBA(client.WriteMultipleCoils(request.StartAddress, request.Quantity, []byte{69}))
			case 16:
				results = errHandlerBA(client.WriteMultipleRegisters(request.StartAddress, request.Quantity, []byte{69}))
			//case 22:
			//	results = errHandlerBA(client.MaskWriteRegister(address, andMask, orMask uint16))
			//case 23:
			//	results = errHandlerBA(client.ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte))
			case 24:
				results = errHandlerBA(client.ReadFIFOQueue(request.StartAddress))
			default:
				continue
			}
			deviceDataChan <- modbusConvertService(results, request.RegisterMap)
		}
	}
}

func errHandlerBA(res []byte, err error) []byte {
	if err != nil {
		log.Printf("RequestsService err %v\n", err)
	}
	return res
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
