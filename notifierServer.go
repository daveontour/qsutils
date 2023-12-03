package qsutils

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"encoding/gob"
)

var listenerMap map[int]chan NotificationServiceMessage
var backlogMap map[int]*List
var clientPropertyMap map[int]map[string]any
var refreshTimerMap map[int]*time.Timer
var clientRegisteredChan chan NotificationClient = make(chan NotificationClient)

// type NotificationServiceMessage struct {
// 	Message     []byte
// 	MessageType int
// }

// type NotificationClient struct {
// 	ProcessID  int
// 	Properties map[string]any
// }

// const (
// 	REGISTERED         = iota
// 	DISABLED           = iota
// 	MESSAGE            = iota
// 	DISCONNECTED       = iota
// 	CLEARBACKLOG       = iota
// 	BROADCASTMESSAGE   = iota
// 	OPERATIONALMESSAGE = iota
// 	TIMEOUT            = iota
// 	REFRESHTIMER       = iota
// )

type NotificationService struct {
	disabled bool // if true, then no messages will be sent to the client
}

func (t *NotificationService) Disable() {
	t.disabled = true
}

func (t *NotificationService) Enable() {
	t.disabled = false
}

func (t *NotificationService) Listen(client NotificationClient, reply *NotificationServiceMessage) error {

	if t.disabled {
		reply.Message = getGobFromString("Disabled")
		reply.MessageType = DISABLED
		return nil
	}

	var dispatcherChan chan NotificationServiceMessage
	var clientProcessID int = client.ProcessID

	//if the channel exists, then use it else create a new channel for the args and use it
	if _, ok := listenerMap[clientProcessID]; ok {
		dispatcherChan = listenerMap[clientProcessID]
	} else {
		dispatcherChan = make(chan NotificationServiceMessage)
	}

	go registerListner(dispatcherChan, client)

	clientRegisteredChan <- client

	message := <-dispatcherChan
	*reply = message

	return nil
}
func (t *NotificationService) ClearBacklog(clientProcessID int, reply *NotificationServiceMessage) error {

	delete(backlogMap, clientProcessID)

	reply.Message = getGobFromString("Backlog Cleared")
	reply.MessageType = CLEARBACKLOG

	return nil
}
func (t *NotificationService) Disconnect(clientProcessID int, reply *NotificationServiceMessage) error {

	delete(listenerMap, clientProcessID)
	delete(backlogMap, clientProcessID)
	delete(clientPropertyMap, clientProcessID)

	if timer, ok := refreshTimerMap[clientProcessID]; ok {
		timer.Stop()
	}
	delete(refreshTimerMap, clientProcessID)

	reply.Message = getGobFromString("Disconnected")
	reply.MessageType = DISCONNECTED

	return nil
}

type registrationHandler func(ch chan NotificationClient)

func InitServer(protocol string, endpoint string, registrationHandler registrationHandler) {

	gob.Register(NotificationServiceMessage{})

	ns := new(NotificationService) // create a new instance of the service
	ns.disabled = false            // set the disabled flag to false

	gob.Register(NotificationServiceMessage{})

	rpc.Register(ns)
	rpc.HandleHTTP()

	go registrationHandler(clientRegisteredChan)

	l, err := net.Listen(protocol, endpoint)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)
}

func registerListner(listenerChan chan NotificationServiceMessage, client NotificationClient) {

	clientID := client.ProcessID
	// check if the map is create and create it if not then create it
	if listenerMap == nil {
		listenerMap = make(map[int]chan NotificationServiceMessage)
	}
	if backlogMap == nil {
		backlogMap = make(map[int]*List)
	}
	if refreshTimerMap == nil {
		refreshTimerMap = make(map[int]*time.Timer)
	}
	if clientPropertyMap == nil {
		clientPropertyMap = make(map[int]map[string]any)
	}

	// if the timer for the client exists, then stop it
	if timer, ok := refreshTimerMap[clientID]; ok {
		timer.Stop()
	}

	delete(refreshTimerMap, clientID)
	delete(listenerMap, clientID)
	delete(clientPropertyMap, clientID)

	listenerMap[clientID] = listenerChan
	clientPropertyMap[clientID] = client.Properties

	if _, ok := backlogMap[clientID]; !ok {
		backlogMap[clientID] = NewList()
	}

	//Set up a timer to send a refresh message to the client every 13 seconds
	clientRefreshTimer := time.NewTimer(13 * time.Second)
	refreshTimerMap[clientID] = clientRefreshTimer

	//create a go routine to wait for the timer to expire and then send a refresh message to the client
	go func(clientID int) {
		<-clientRefreshTimer.C
		m1 := NotificationServiceMessage{Message: getGobFromString("Refresh Timer"), MessageType: REFRESHTIMER}
		go sendMessageToClient(clientID, m1)
		delete(refreshTimerMap, clientID)
	}(clientID)

	//check if there are any messages in the backlog and send them to the client
	processBacklog(clientID)
}

func processBacklog(clientID int) {
	if backlogMessage, hasBacklog := backlogMap[clientID].FrontPop(); hasBacklog {
		sendMessageToClient(clientID, backlogMessage.Value.(NotificationServiceMessage))
	}
}

func SendMessageToClient(clientID int, message NotificationServiceMessage) {
	sendMessageToClient(clientID, message)
}

func SendBroadcastMessage(message NotificationServiceMessage) {
	for clientID := range listenerMap {
		sendMessageToClient(clientID, message)
	}
}

func SendMessageToClients(clientIDs []int, message NotificationServiceMessage) {
	for _, clientID := range clientIDs {
		sendMessageToClient(clientID, message)
	}
}

func sendMessageToClient(clientID int, message NotificationServiceMessage) {

	if l, ok := listenerMap[clientID]; ok {
		delete(listenerMap, clientID)
		l <- message
	} else {
		backlogMap[clientID].PushBack(message)
	}
}
