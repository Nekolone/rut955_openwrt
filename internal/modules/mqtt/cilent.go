package mqtt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Clients struct {
	Clients []Client `json:"clients"`
}

type Client struct {
	Name             string     `json:"name"`
	IPPort           string     `json:"ip_port"`
	SubscriptionList []string   `json:"subscription_list"`
	DataFormat       DataFormat `json:"data_format"`
}

type DataFormat struct {
	DataNameField  string `json:"data_name_field"`
	DataTypeField  string `json:"data_type_field"`
	DataValueField string `json:"data_value_field"`
}

func Start(dataSourceChan chan map[string][]string, path string) {
	defer log.Print("mqtt ds - done")
	log.Print("connect to mqtt data source")

	mqttConfig := getMqttConfig(path)

	for _, clientConfig := range mqttConfig.Clients {
		go clientHandler(clientConfig, dataSourceChan)
	}
}

func clientHandler(clientConfig Client, dataSourceChan chan map[string][]string) {
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

	log.Print("mqtt module successfully start")
	for range time.NewTicker(time.Second * 10).C {
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			continue
		}
		SubscribeService(mqttClient, subMap, dataSourceChan, clientConfig.DataFormat, clientConfig.Name)
		log.Print("reconnect to mqtt broker")
	}
}

func SubscribeService(c MQTT.Client, subMap map[string]byte, dataSourceChan chan map[string][]string, dataFormat DataFormat, moduleName string) {
	defer c.Disconnect(250)

	if token := c.SubscribeMultiple(subMap, CallbackModifier(dataSourceChan, dataFormat, moduleName)); token.Wait() &&
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

func CallbackModifier(dataSourceChan chan map[string][]string, dataFormat DataFormat, moduleName string) MQTT.MessageHandler {
	var res MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		dataSourceChan <- getValueFromJSON(dataFormat, msg.Payload(), moduleName)
	}
	return res
}

func getValueFromJSON(dataFormat DataFormat, payload []byte, name string) map[string][]string {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Getting data from json error. Error > %v\nMSG>%s", r, payload)
		}
	}()
	data := make(map[string]interface{})
	if err := json.Unmarshal(payload, &data); err != nil {
		return map[string][]string{}
	}
	res := make(map[string]interface{})
	Flatten2("", data, res)
	dataType := res[dataFormat.DataTypeField]
	if res[dataFormat.DataTypeField] == "" || res[dataFormat.DataTypeField] == nil {
		dataType = "3"
	}

	return map[string][]string{
		name: {
			getCurTime(),
			fmt.Sprintf(
				"%v|%s|%v",
				res[dataFormat.DataNameField],
				dataType,
				res[dataFormat.DataValueField],
			),
		},
	}

	// return fmt.Sprintf(
	// 	"%s:%s:%s",
	// 	res[dataFormat.DataNameField],
	// 	dataType,
	// 	res[dataFormat.DataValueField],
	// )
}

func getCurTime() string {
	out, err := exec.Command("gpsctl", "-e").Output()
	if err != nil || bytes.Equal(out, []byte("1970-01-01 02:00:00")) {
		out = []byte(time.Now().Format("2006-01-02 15:04:05"))
	}
	return string(out[8:10]) + string(out[5:7]) + string(out[2:4]) + string(out[11:13]) + string(out[14:16]) + string(out[17:19])
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

func setDefaultMqttConfig() *Clients {
	return &Clients{
		[]Client{{
			Name:             "unknown device",
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
