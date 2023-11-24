package main

import (
	"fmt"
	"time"

	qsutils "github.com/daveontour/qsutils"
)

func main() {
	// Create a new list
	list := qsutils.New()

	// Add some elements to the list

	// list.PushFront(1)
	// list.PushFront(2)
	// list.PushFront(3)

	list.PushFront("No Date")
	list.InsertByDateTime("Now", time.Now())
	list.InsertByDateTime("Latest", time.Now().Add(time.Hour*2))
	list.InsertByDateTime("Earliest", time.Now().Add(time.Hour*-2))
	list.InsertByDateTime("Later", time.Now().Add(time.Hour*1))
	list.PushBack("Also No Date")

	//iterate through the list and print the values

	for e := list.FrontOnly(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

}
