package main

import (
	"fmt"

	qsutils "github.com/daveontour/qsutils"
)

func main() {

	go qsutils.InitServer("tcp", ":1234", registrationHandler)
	go qsutils.ListenExample("tcp", "localhost:1234")

	ch := make(chan int)
	<-ch
}

func registrationHandler(clientRegisteredChan chan qsutils.NotificationClient) {
	for {
		m := <-clientRegisteredChan
		fmt.Println("Client registered: ", m)
	}
}
