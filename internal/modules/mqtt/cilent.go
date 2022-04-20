package mqtt

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"strconv"
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
	defer log.Print("mqtt ds - done")
	log.Print("connect to mqtt data source")

	mqttConfig := getMqttConfig(path)

	for _, clientConfig := range mqttConfig.Clients {
		go clientHandler(clientConfig, dataSourceChan)
	}
}

func clientHandler(clientConfig Client, dataSourceChan chan string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("something goes wrong in mqtt module.\nclientConfig - %v\nMsg > %v", clientConfig, r)
		}
	}()
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

func SubscribeService(c MQTT.Client, subMap map[string]byte, dataSourceChan chan string, dataFormat DataFormat) {
	defer c.Disconnect(250)

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

func CallbackModifier(dataSourceChan chan string, dataFormat DataFormat) MQTT.MessageHandler {
	var res MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		dataSourceChan <- getValueFromJSON(dataFormat, msg.Payload())
	}
	return res
}

func getValueFromJSON(dataFormat DataFormat, payload []byte) string {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Getting data from json error. Error > %v\nMSG>%s", r, payload)
		}
	}()
	data := make(map[string]interface{})
	if err := json.Unmarshal(payload, &data); err != nil {
		return ""
	}
	res := make(map[string]interface{})
	Flatten2("", data, res)
	if res[dataFormat.DataTypeField] == "" {
		res[dataFormat.DataTypeField] = 3
	}

	return fmt.Sprintf(
		"%s:%s:%s",
		res[dataFormat.DataNameField],
		res[dataFormat.DataTypeField],
		res[dataFormat.DataValueField],
	)
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

func getMqttConfig(path string) (cfg *Clients) {
	cfg = setDefaultMqttConfig()
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

func setDefaultMqttConfig() *Clients {
	return &Clients{
		[]Client{{
			IPPort:           "127.0.0.1:18883",
			SubscriptionList: []string{"#"},
			DataFormat: DataFormat{
				DataNameField:  "req_name",
				DataTypeField:  "type",
				DataValueField: "data_value",
			},
		}},
	}
}
