package mqtt

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/TheThingsNetwork/ttn/core/types"
	//. "github.com/TheThingsNetwork/ttn/utils/testing"
	"github.com/apex/log"
	//. "github.com/smartystreets/assertions"
)

// ----- MAIN -----
func main() {
	connectAndSubscribe()
}

// ----- PAYLOAD DECLARATION -----
func definePayload(req types.UplinkMessage) io.Reader {

	type Payload struct {
		Data string `json:"dataToSend"`
	}

	data := Payload{
		// fill struct depending on ttn msg
		"data test",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		print("ERROR while JSON encoding ")
	}
	body := bytes.NewReader(payloadBytes)
	return body
}

// ----- SENDING TO OPENSENSORS -----
func sendMsgToOS(body io.Reader) {
	req, err := http.NewRequest("POST", "https://realtime.opensensors.io/v1/topics//users/aroy314/MyFirstDevice?client-id=5596&password=yuFZtRqi", body)
	if err != nil {
		print("ERROR NewRequest : " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "api-key 3deb3232-95ce-43be-8933-7e52145d48c0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		print("ERROR request exec : " + err.Error())
	}
	defer resp.Body.Close()
}

// ----- CONNECTION, SUBSCRIPTION -----
func connectAndSubscribe() {
	//CONNECTION SETUP
	ctx := log.WithField("MyExample", "Go Client")
	client := NewClient(
		ctx,
		"ttnctl",
		"office-app",
		"ttn-account-preview.OfuuW9smtu33PjpPtVAs54Bmc2dcgHEOywtuAT1oqzk",
		"tcp://eu.thethings.network:1883",
	)
	//CONNECTION ERROR
	if err := client.Connect(); err != nil {
		ctx.WithError(err).Fatal("ERROR : Could not connect")
	}

	// SUBSCRIBE TO DEVICE
	token := client.SubscribeDeviceUplink("office-app", "dev-id", func(client Client, appID string, devID string, req types.UplinkMessage) {
		//get message structure
		print("message analysis... ")
		message := definePayload(req)
		// Do something with the uplink message
		print("sending... ")
		sendMsgToOS(message)
		print("sent ! ")
	})
	token.Wait()
	if err := token.Error(); err != nil {
		ctx.WithError(err).Fatal("ERROR : Could not subscribe")
	}
}

//working "hello world" curl command :
//curl -X POST --header "Content-Type: application/json" --header "Authorization: api-key 3deb3232-95ce-43be-8933-7e52145d48c0" -d '{"data": "Hello World"}' "https://realtime.opensensors.io/v1/topics//users/aroy314/MyFirstDevice?client-id=5596&password=yuFZtRqi"

//with mosquitto :
//mosquitto_sub -h eu.thethings.network:1883 -d -t 'my-app-id/devices/my-dev-id/up'
//mosquitto_pub -h mqtt.opensensors.io -i 5596 -t /users/aroy314/test -u aroy314 -m 'My first message' -P yuFZtRqi
