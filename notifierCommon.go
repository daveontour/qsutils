package qsutils

import (
	"bytes"
	"encoding/gob"
)

type NotificationServiceMessage struct {
	Message     []byte
	MessageType int
}

type NotificationClient struct {
	ProcessID  int
	Properties map[string]any
}

const (
	REGISTERED         = iota
	DISABLED           = iota
	MESSAGE            = iota
	DISCONNECTED       = iota
	CLEARBACKLOG       = iota
	BROADCASTMESSAGE   = iota
	OPERATIONALMESSAGE = iota
	TIMEOUT            = iota
	REFRESHTIMER       = iota
)

func getStringFromGob(message []byte) string {
	var str string
	err := gob.NewDecoder(bytes.NewReader(message)).Decode(&str)
	if err != nil {
		return err.Error()
	}

	return str
}
func getGobFromString(message string) []byte {
	var bBuf bytes.Buffer
	gob.NewEncoder(&bBuf).Encode(&message)

	return bBuf.Bytes()
}
