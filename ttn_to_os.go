package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	classic_log "log"
	"net/http"
	"os"
	"sort"

	"github.com/TheThingsNetwork/ttn/core/types"
	"github.com/TheThingsNetwork/ttn/mqtt"
	"github.com/apex/log"
	"github.com/callum-ramage/jsonconfig" //conf json lib
)

type Configuration struct {
	ttnAppID  string `json:"ttnAppID"`
	ttnAPIkey string `json:"ttnAPIkey"`
	osURL     string `json:"osURL"`
	osAPIkey  string `json:"osAPIkey"`
}

// ----- MAIN -----
func main() {

	//get config from json file
	config := getConfig("config.json")

	//launch the client with the good config
	connectAndSubscribe(config)
}

// ----- CONNECTION, SUBSCRIPTION -----
func connectAndSubscribe(conf Configuration) {
	//CONNECTION SETUP
	ctx := log.WithField("MyExample", "Go Client")
	client := mqtt.NewClient(
		nil,
		"ttnctl",
		conf.ttnAppID,
		conf.ttnAPIkey,
		"tcp://eu.thethings.network:1883",
	)

	//CONNECTION ERROR
	if err := client.Connect(); err != nil {
		ctx.WithError(err).Fatal("ERROR : Could not connect")
	}

	// SUBSCRIBE TO DEVICE
	token := client.SubscribeAppUplink(conf.ttnAppID, func(client mqtt.Client, appID string, devID string, req types.UplinkMessage) {

		print("Incoming message : " + appID + " " + devID + " \n")
		setPayloadAndSend(conf, req)
		print("Sent ! \n")

	})
	token.Wait()
	if err := token.Error(); err != nil {
		ctx.WithError(err).Fatal("ERROR : Could not subscribe")
	}

	//wait for user input
	print("Press 'enter' to exit \n")
	waitKeyPressed()

}

// ----- PAYLOAD DECLARATION -----
func setPayloadAndSend(conf Configuration, msg types.UplinkMessage) {

	type Payload struct {
		Data string `json:"data"`
	}

	//displaying fields sent
	println("payloadFields : ")
	var keys []string
	for k := range msg.PayloadFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, " : ", msg.PayloadFields[k])
	}

	//put the map in json format
	payloadBytes, err := json.Marshal(msg.PayloadFields)
	if err != nil {
		print("ERROR while JSON encoding ")
	}

	//put this JSON in the Payload Data field
	data := Payload{
		string(payloadBytes),
	}

	//put the Payload in JSON format
	dataBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(dataBytes)

	//set the request
	req, err := http.NewRequest("POST", conf.osURL, body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "api-key "+conf.osAPIkey)

	//send it !
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

}

func getConfig(filename string) Configuration {

	config, err := jsonconfig.LoadAbstract(filename, "")
	if err != nil {
		fmt.Println("error loadabstract:", err)
	}

	conf := Configuration{
		ttnAppID:  config["ttnAppID"].Str,
		ttnAPIkey: config["ttnAPIkey"].Str,
		osURL:     config["osURL"].Str,
		osAPIkey:  config["osAPIkey"].Str,
	}

	return conf
}

func waitKeyPressed() bool {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		classic_log.Fatal(err)
	}
	return response == "\n"
}

//working "hello world" curl command :
//curl -X POST --header "Content-Type: application/json" --header "Authorization: api-key 3deb3232-95ce-43be-8933-7e52145d48c0" -d '{"data": "Hello World"}' "https://realtime.opensensors.io/v1/topics//users/aroy314/MyFirstDevice?client-id=5596&password=yuFZtRqi"

//with mosquitto :
//mosquitto_sub -h eu.thethings.network:1883 -d -t 'my-app-id/devices/my-dev-id/up'
//mosquitto_pub -h mqtt.opensensors.io -i 5596 -t /users/aroy314/test -u aroy314 -m 'My first message' -P yuFZtRqi
