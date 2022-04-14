package mqtt

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"strings"
	"time"
)

type Clients struct {
	Clients []Client `json:"clients"`
}

type Client struct {
	IPPort           string     `json:"ip_port"`
	SubscriptionList []string   `json:"subscription_list"`
	DataFormat       DataFormat `json:"data_format"`
}

type DataFormat struct {
	DataNameField  string `json:"data_name_field"`
	DataTypeField  string `json:"data_type_field"`
	DataValueField string `json:"data_value_field"`
}

func Start(dataSourceChan chan string, path string) {
	log.Println()
	mqttConfig, err := getMqttConfig(path)
	if err != nil {
		log.Println("mqtt config path error, default config for local mqtt broker")
		mqttConfig = setDefaultMqttConfig()
	}

	for _, clientConfig := range mqttConfig.Clients {
		go clientHandler(clientConfig, dataSourceChan)
	}
}

func CallbackModifier(dataSourceChan chan string, dataFormat DataFormat) MQTT.MessageHandler {
	var res MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		dataSourceChan <- getValueFromJSON(dataFormat, msg.Payload())
	}
	return res
}

func getValueFromJSON(dataFormat DataFormat, payload []byte) string {
	jsonMap := make(map[string]json.RawMessage)
	err := json.Unmarshal(payload, &jsonMap)
	if err != nil {
		return ""
	}
	payloadMap := make(map[string]string)
	for k, message := range jsonMap {
		payloadMap[k] = string(message)
	}
	if payloadMap[dataFormat.DataTypeField] == "" {
		payloadMap[dataFormat.DataTypeField] = "3"
	}

	return fmt.Sprintf(
		"%s:%s:%s",
		strings.ReplaceAll(payloadMap[dataFormat.DataNameField], "\"", ""),
		strings.ReplaceAll(payloadMap[dataFormat.DataTypeField], "\"", ""),
		payloadMap[dataFormat.DataValueField],
	)
	// return map[string]string{
	//	"name":  payloadMap[dataFormat.DataNameField],
	//	"type":  payloadMap[dataFormat.DataTypeField],
	//	"value": payloadMap[dataFormat.DataValueField],
	//}
}

func clientHandler(clientConfig Client, dataSourceChan chan string) {
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", clientConfig.IPPort))
	mqttClient := MQTT.NewClient(opts)

	subMap := make(map[string]byte)
	for _, s := range clientConfig.SubscriptionList {
		subMap[s] = 0
	}

	for range time.NewTicker(time.Second * 10).C {
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			continue
		}
		SubscribeService(mqttClient, subMap, dataSourceChan, clientConfig.DataFormat)
	}
}

//
// func convertMqttData(deviceDataChan chan string) {
//	deviceDataChan <- <-subCallBackChan
//}

func SubscribeService(c MQTT.Client, subMap map[string]byte, dataSourceChan chan string, dataFormat DataFormat) {
	defer c.Disconnect(250)

	// if token := c.SubscribeMultiple(subMap, msgHandl); token.Wait() && token.Error() != nil {
	// if token := c.SubscribeMultiple(map[string]byte{"modbusData": 0, "modbusData2": 0}, CallbackModifier(deviceDataChan,
	if token := c.SubscribeMultiple(subMap, CallbackModifier(dataSourceChan, dataFormat)); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		return
	}

	defer c.Unsubscribe("#")

	for range time.NewTicker(time.Second * 10).C {
		if c.IsConnected() {
			continue
		}
		return
	}
}

func setDefaultMqttConfig() *Clients {
	return &Clients{
		[]Client{{
			IPPort:           "127.0.0.1:18883",
			SubscriptionList: []string{"#"},
			DataFormat: DataFormat{
				DataNameField:  "name",
				DataTypeField:  "type",
				DataValueField: "value",
			},
		}},
	}
}

func getMqttConfig(path string) (*Clients, error) {
	var cfg Clients
	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return &Clients{}, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	_ = jsonParser.Decode(&cfg)
	return &cfg, nil
}
