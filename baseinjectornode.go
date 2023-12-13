package qsutils

import (
	"strconv"
)

// SourceDestination is an interface that implements a method to get the source and destination of a message
type SourceDestination interface {
	SetConfig(map[string]string)
	Prepare(chan string)
	StartListening(chan Datagram) error
	StartSending(chan Datagram)
	Stop()
	ReportStatus() string
}
type BaseNode struct {
	reportChan   chan string
	MessagesSent int
	TriggerID    string
	Execute      bool
	nodeconfig   map[string]string
	Interval     int
	MaxRuntime   int
	MaxMessages  int
	MinInterval  int
	MaxInterval  int
	SendType     int
	Source       string
	Destination  string
	Enabled      bool
}

func (r *BaseNode) SetConfig(config map[string]string) {
	r.nodeconfig = config
	r.setConfig()
}

func (r *BaseNode) setConfig() {
	r.TriggerID = r.nodeconfig["TriggerID"]
	r.Source = r.nodeconfig["Source"]
	r.Destination = r.nodeconfig["Destination"]
	r.Interval, _ = strconv.Atoi(r.nodeconfig["Interval"])
	r.MaxRuntime, _ = strconv.Atoi(r.nodeconfig["MaxRuntime"])
	r.MaxMessages, _ = strconv.Atoi(r.nodeconfig["MaxMessages"])
	r.MinInterval, _ = strconv.Atoi(r.nodeconfig["MinimumInterval"])
	r.MaxInterval, _ = strconv.Atoi(r.nodeconfig["MaximumInterval"])

	if r.nodeconfig["SendType"] == "BURSTSEND" {
		r.SendType = BURSTSEND
	} else if r.nodeconfig["SendType"] == "RANDOMSEND" {
		r.SendType = RANDOMSEND
	} else {
		r.SendType = NORMALSEND
	}

	r.Enabled, _ = strconv.ParseBool(r.nodeconfig["Enabled"])
}

func (r *BaseNode) setIntegerConfigValue(key string, value int) {
	if r.nodeconfig == nil {
		r.nodeconfig = make(map[string]string)
	}
	r.nodeconfig[key] = strconv.Itoa(value)
	r.setConfig()
}
func (r *BaseNode) getIntegerConfigValue(key string) int {
	if i, ok := r.nodeconfig["interval"]; ok {
		result, _ := strconv.Atoi(i)
		return result
	} else {
		return 0
	}
}

func (r *BaseNode) setStringConfigValue(key string, value string) {
	if r.nodeconfig == nil {
		r.nodeconfig = make(map[string]string)
	}
	r.nodeconfig[key] = value
	r.setConfig()
}

func (r *BaseNode) getStringConfigValue(key string) string {
	if i, ok := r.nodeconfig[key]; ok {
		return i
	} else {
		return ""
	}
}

func (r *BaseNode) GetBooleanConfigValue(key string) bool {
	if i, ok := r.nodeconfig[key]; ok {
		result, _ := strconv.ParseBool(i)
		return result
	} else {
		return false
	}
}

func (r *BaseNode) setBooleanConfigValue(key string, value bool) {
	if r.nodeconfig == nil {
		r.nodeconfig = make(map[string]string)
	}
	r.nodeconfig[key] = strconv.FormatBool(value)
	r.setConfig()
}

func (r *BaseNode) SetReportChan(reportChan chan string) {
	r.reportChan = reportChan
}

func (r *BaseNode) GetTriggerID() string {
	return r.getStringConfigValue("triggerID")
}

func (r *BaseNode) SetNodeConfig(nodeconfig map[string]string) {
	r.nodeconfig = nodeconfig
}

func (r *BaseNode) SetTriggerID(triggerID string) {
	r.setStringConfigValue("triggerID", triggerID)
}

func (r *BaseNode) SetInterval(interval int) {
	r.setIntegerConfigValue("interval", interval)
}

func (r *BaseNode) SetMaxMessages(maxMessages int) {
	r.setIntegerConfigValue("maxMessages", maxMessages)
}

func (r *BaseNode) SetEnabled() {
	r.setBooleanConfigValue("enabled", true)
}
func (r *BaseNode) SetDisabled() {
	r.setBooleanConfigValue("enabled", false)
}
