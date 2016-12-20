
# ttn_to_opensensors
go project to connect TheThingsNetwork to OpenSensors

## set up
* create your device+topic on opensensors website
* make sure you have go installed
* `go get github.com/TheThingsNetwork/ttn`
* `go get github.com/callum-ramage/jsonconfig`
* create a config.json file from config_sample.json and put your parameters inside
* cd in the folder then `go build ttn_to_os.go` will create the executable
