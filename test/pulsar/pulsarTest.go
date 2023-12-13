package main

import (
	"github.com/daveontour/qsutils"
)

func main() {

	p := qsutils.Pulsar{}
	//create a mapa of string to string

	config := make(map[string]string)
	// config["TriggerID"] = "Pulsar"
	// config["Source"] = "RabbitMQ"
	// config["Destination"] = "RabbitMQ"
	config["Interval"] = "1000"
	// config["MaxRuntime"] = "20000"
	config["MaxMessages"] = "10"
	// config["MinimumInterval"] = "100"
	// config["MaximumInterval"] = "5000"
	// config["Enabled"] = "true"
	config["SendType"] = "NORMALSEND"

	reportChan := make(chan string)
	dataChan := make(chan qsutils.Datagram)
	blockChan := make(chan string)

	go func() {
		for {
			select {
			case report := <-reportChan:
				println(report)
			case data := <-dataChan:
				if data.LastData {
					blockChan <- "done"
					return
				}
				println(string(data.Data))
			}
		}
	}()
	p.SetConfig(config)
	p.Prepare(reportChan)

	go p.StartSending(dataChan)

	<-blockChan
}
