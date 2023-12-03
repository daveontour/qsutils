package qsutils

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"
)

func Listen(protocol, address string) (reply *NotificationServiceMessage) {
	client, err := rpc.DialHTTP(protocol, address)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var clientProcessID int = os.Getpid()
	notificationClient := NotificationClient{ProcessID: clientProcessID}

	// Asynchronous call
	reply = new(NotificationServiceMessage)
	divCall := client.Go("NotificationService.Listen", notificationClient, reply, nil)
	<-divCall.Done // will be equal to divCall
	return
}

// Example of how to use the Listen function
func ListenExample(protocol, address string) {

	for {

		reply := Listen(protocol, address)

		switch reply.MessageType {
		case MESSAGE:
			fmt.Println("Message received: ", getStringFromGob(reply.Message))
		case TIMEOUT:
			fmt.Println("Timeout received: ", getStringFromGob(reply.Message))
		case BROADCASTMESSAGE:
			fmt.Println("Broadcast message received received: ", getStringFromGob(reply.Message))
		case DISABLED:
			fmt.Println("Server currently disabled. Will try again in 10 seconds")
			time.Sleep(10 * time.Second)
		case REFRESHTIMER:
			fmt.Println("Refresh timer received: ", getStringFromGob(reply.Message))

		default:
			fmt.Println("Notification received: ", getStringFromGob(reply.Message))
		}
	}
}
