package qsutils

import (
	"errors"
	"math/rand"
	"time"
)

type SourceDestinationConfig struct {
}
type Datagram struct {
	TriggerID string
	Data      []byte
}

// SourceDestination is an interface that implements a method to get the source and destination of a message
type SourceDestination interface {
	GetSource() string
	GetDestinationDescription() string
	SetConfig(config SourceDestinationConfig)
	Prepare()
	StartListening(chan []byte) error
	StartSending(chan []byte)
	Stop()
	ReportStatus() string
}

type BaseNode struct {
	reportChan              chan string
	execute                 bool
	triggerID               string
	interval                int
	minimumInterval         int
	maximumInterval         int
	maxMessages             int
	maxRuntime              int
	repeatExecution         bool
	repeatExecutionInterval int
}

type Pulsar struct {
	BaseNode
	config                 SourceDestinationConfig
	source                 string
	destinationDescription string
	burstSend              bool
	regularSend            bool
	randomSend             bool
}

func (r *Pulsar) GetSource() string {
	return r.source
}

func (r *Pulsar) GetDestinationDescription() string {
	return r.destinationDescription
}

func (r *Pulsar) SetConfig(config SourceDestinationConfig) {
	r.config = config
}

func (r *Pulsar) Prepare(reportChan chan string) {
	r.reportChan = reportChan
	r.source = "RabbitMQ"
	r.destinationDescription = "RabbitMQ"
}

func (r *Pulsar) StartListening(ch chan Datagram) error {
	return errors.New("not implemented")
}

func (r *Pulsar) Stop() {
	r.execute = false
}

func (r *Pulsar) getDatagram() Datagram {
	return Datagram{TriggerID: r.triggerID, Data: []byte("Pulse")}
}
func (r *Pulsar) StartSending(ch chan Datagram) {

	r.execute = true
	datagram := r.getDatagram()
	messagesSent := 0
	interval := r.interval

	if r.maxRuntime > 0 {
		go func() {
			time.Sleep(time.Duration(r.maxRuntime) * time.Second)
			r.execute = false
		}()
	}

	for r.execute && (r.maxMessages > messagesSent || r.maxMessages < 0) {

		ch <- datagram
		messagesSent++

		if r.randomSend {
			interval = rand.Intn(r.maximumInterval-r.minimumInterval) + r.minimumInterval
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
