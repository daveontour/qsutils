package qsutils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Datagram struct {
	TriggerID string
	Data      []byte
	LastData  bool
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Compiler will check that Pulsar implements SourceDestination
var _ SourceDestination = new(Pulsar)
var _ SourceDestination = new(RabbitMQNode)

const (
	BURSTSEND = iota
	RANDOMSEND
	NORMALSEND
	EVENTSEND
)

type Pulsar struct {
	BaseNode
}

func (r *Pulsar) ReportStatus() string {
	return ""
}
func (r *Pulsar) SetConfig(config map[string]string) {
	r.nodeconfig = config
	r.setConfig()
}
func (r *Pulsar) Prepare(reportChan chan string) {
	r.reportChan = reportChan
}
func (r *Pulsar) StartListening(ch chan Datagram) error {
	return errors.New("not implemented")
}
func (r *Pulsar) Stop() {
	r.Execute = false
}
func (r *Pulsar) getDatagram() Datagram {
	return Datagram{TriggerID: r.TriggerID, Data: []byte(fmt.Sprintf("Pulse Number %v", r.MessagesSent))}
}
func (r *Pulsar) StartSending(ch chan Datagram) {

	r.Execute = true
	datagram := r.getDatagram()
	r.MessagesSent = 0
	interval := r.Interval

	if r.MaxRuntime > 0 {
		go func() {
			time.Sleep(time.Duration(r.MaxRuntime) * time.Second)
			r.Execute = false
		}()
	}

	for r.Execute && (r.MaxMessages > r.MessagesSent || r.MaxMessages <= 0) {

		ch <- datagram
		r.MessagesSent++
		datagram = r.getDatagram()

		if r.SendType == BURSTSEND {
			continue
		}

		if r.SendType == RANDOMSEND {
			interval = rand.Intn(r.MaxInterval-r.MinInterval) + r.MinInterval
		}

		if r.Execute && (r.MaxMessages > r.MessagesSent || r.MaxMessages <= 0) {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}

	}

	r.reportChan <- fmt.Sprintf("Pulsar %v sent %v messages", r.TriggerID, r.MessagesSent)

	datagram.LastData = true
	ch <- datagram

	close(ch)
	close(r.reportChan)
}

type RabbitMQNode struct {
	BaseNode
	RabbitMQConnectionString string
}

func (r *RabbitMQNode) ReportStatus() string {
	return ""
}

func (r *RabbitMQNode) SetConfig(config map[string]string) {
	r.nodeconfig = config
	r.setConfig()
}

func (r *RabbitMQNode) Prepare(reportChan chan string) {
	r.reportChan = reportChan
}

func (r *RabbitMQNode) StartListening(ch chan Datagram) error {
	return errors.New("not implemented")
}

func (r *RabbitMQNode) Stop() {
	r.Execute = false
}

func (r *RabbitMQNode) getDatagram() Datagram {
	return Datagram{TriggerID: r.TriggerID, Data: []byte(fmt.Sprintf("RabbitMQ Message %v", r.MessagesSent))}
}

func (r *RabbitMQNode) StartSending(ch chan Datagram) {

	conn, err := amqp.Dial(r.RabbitMQConnectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	rmqCh, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer rmqCh.Close()

	r.Execute = true
	datagram := r.getDatagram()
	r.MessagesSent = 0
	interval := r.Interval

	if r.MaxRuntime > 0 {
		go func() {
			time.Sleep(time.Duration(r.MaxRuntime) * time.Second)
			r.Execute = false
		}()
	}

	for r.Execute && (r.MaxMessages > r.MessagesSent || r.MaxMessages <= 0) {
		ch <- datagram
		r.MessagesSent++
		datagram = r.getDatagram()

		if r.SendType == BURSTSEND || r.SendType == EVENTSEND {
			continue
		}

		if r.SendType == RANDOMSEND {
			interval = rand.Intn(r.MaxInterval-r.MinInterval) + r.MinInterval
		}

		if r.Execute && (r.MaxMessages > r.MessagesSent || r.MaxMessages <= 0) {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}

	r.reportChan <- fmt.Sprintf("RabbitMQ Node %v sent %v messages", r.TriggerID, r.MessagesSent)

	datagram.LastData = true
	ch <- datagram

	close(ch)
	close(r.reportChan)
}

func ReadCSV() {

	// os.Open() opens specific file in
	// read-only mode and this return
	// a pointer of type os.File
	file, err := os.Open("Students.csv")

	// Checks for the error
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	// Closes the file
	defer file.Close()

	// The csv.NewReader() function is called in
	// which the object os.File passed as its parameter
	// and this creates a new csv.Reader that reads
	// from the file
	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file
	// and Returns them as slice of slices of string
	// and an error if any
	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		fmt.Println("Error reading records")
	}

	// Loop to iterate through
	// and print each of the string slice
	for _, eachrecord := range records {
		fmt.Println(eachrecord)
	}
}

// A function to write to a CSV file
func WriteCSV() {
	// Creating a file
	file, err := os.Create("Students.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	// Creating a writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Writing to the file
	student1 := []string{"John", "Doe", "CS"}
	student2 := []string{"Mary", "Moe", "IT"}
	student3 := []string{"Jane", "Doe", "IT"}
	student4 := []string{"Mike", "Tee", "CS"}
	student5 := []string{"Kate", "Lee", "CS"}

	// Writing individual student to the file
	err = writer.Write(student1)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}
	err = writer.Write(student2)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}
	err = writer.Write(student3)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}
	err = writer.Write(student4)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}
	err = writer.Write(student5)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}

	// Writing multiple students to the file
	students := [][]string{
		{"John", "Doe", "CS"},
		{"Mary", "Moe", "IT"},
		{"Jane", "Doe", "IT"},
		{"Mike", "Tee", "CS"},
		{"Kate", "Lee", "CS"},
	}

	err = writer.WriteAll(students)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}
}
