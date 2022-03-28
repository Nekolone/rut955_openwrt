package mqtt

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)

type Client struct {
	IpPort           string     `json:"ip_port"`
	SubscriptionList []string   `json:"subscription_list"`
	DataFormat       DataFormat `json:"data_format"`
}

type DataFormat struct {
	DataNameField  string `json:"data_name_field"`
	DataTypeField  string `json:"data_type_field"`
	DataValueField string `json:"data_value_field"`
}

type ClientList struct {
	Clients []Client `json:"client_list"`
}

func Start(deviceDataChan chan string, path string) {
	clientList, err := getMqttConfig(path)
	if err != nil {
		log.Println("mqtt config path error, default config for local mqtt broker")
		clientList = setDefaultMqttConfig()
	}

	for _, client := range clientList.Clients {
		go clientHandler(&client, deviceDataChan)
	}
}

func CallbackModifier(deviceDataChan chan string, dataFormat DataFormat) MQTT.MessageHandler {
	var res MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		deviceDataChan <- getValueFromJson(dataFormat, msg.Payload())
	}
	return res
}

func getValueFromJson(dataFormat DataFormat, payload []byte) string {
	jsonMap := make(map[string]json.RawMessage)
	err := json.Unmarshal(payload, &jsonMap)
	if err != nil {
		return ""
	}
	payloadMap := make(map[string]string)
	for k, message := range jsonMap {
		payloadMap[k] = string(message)
	}
	return fmt.Sprintf("%s:%s:%s", payloadMap[dataFormat.DataNameField], payloadMap[dataFormat.DataTypeField],
		payloadMap[dataFormat.DataValueField])
}

func clientHandler(client *Client, deviceDataChan chan string) {
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", client.IpPort))
	c := MQTT.NewClient(opts)
	subMap := make(map[string]byte)

	for _, s := range client.SubscriptionList {
		subMap[s] = 0
	}

	for range time.NewTicker(time.Second * 10).C {
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			continue
		}
		SubscribeService(c, subMap, deviceDataChan, client.DataFormat)

	}

}

//
//func convertMqttData(deviceDataChan chan string) {
//	deviceDataChan <- <-subCallBackChan
//}

func SubscribeService(c MQTT.Client, subMap map[string]byte, deviceDataChan chan string, dataFormat DataFormat) {
	defer c.Disconnect(250)

	//if token := c.SubscribeMultiple(subMap, msgHandl); token.Wait() && token.Error() != nil {
	//if token := c.SubscribeMultiple(map[string]byte{"modbusData": 0, "modbusData2": 0}, CallbackModifier(deviceDataChan,
	if token := c.SubscribeMultiple(subMap, CallbackModifier(deviceDataChan, dataFormat))
		token.Wait() && token.Error() != nil {

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

func setDefaultMqttConfig() ClientList {
	return ClientList{
		[]Client{{
			IpPort:           "127.0.0.1:1883",
			SubscriptionList: []string{"#"},
			DataFormat: DataFormat{
				DataNameField:  "name",
				DataTypeField:  "type",
				DataValueField: "value",
			},
		}},
	}
}

func getMqttConfig(path string) (cfg ClientList, err error) {
	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return ClientList{}, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	_ = jsonParser.Decode(&cfg)
	return cfg, nil
}
