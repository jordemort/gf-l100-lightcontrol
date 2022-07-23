package main

import (
	"fmt"
	"net/http"
	"os"

	"go.bug.st/serial"
)

const portName = "/dev/ttyAMA1"
const initCmd = "0xfa1"
const motionCmd = "0xb6a"
const onCmd = "0xbb9"
const offCmd = "0xbba"

type kkSwitch struct {
	name      string
	state     bool
	statePath string
}

var lightSwitch = kkSwitch{
	name:      "light",
	state:     false,
	statePath: "/mnt/mtd/ipc/light.state",
}

var motionSwitch = kkSwitch{
	name:      "motion",
	state:     false,
	statePath: "/mnt/mtd/ipc/motion.state",
}

var relay serial.Port

func sendCmd(cmd string) error {
	n, err := relay.Write([]byte(cmd))
	if err != nil {
		return err
	}
	if n != len(cmd) {
		return fmt.Errorf("only wrote %d of %d bytes", n, len(cmd))
	}
	return nil
}

func makeState(isOn bool) string {
	if isOn {
		return "on"
	} else {
		return "off"
	}
}

func updateRelay() error {
	if lightSwitch.state {
		return sendCmd(onCmd)
	} else if motionSwitch.state {
		return sendCmd(motionCmd)
	} else {
		return sendCmd(offCmd)
	}
}

func (s *kkSwitch) writeState() error {
	data := []byte(makeState(s.state))
	return os.WriteFile(s.statePath, data, 0644)
}

func (s *kkSwitch) initState() {
	rawState, err := os.ReadFile(s.statePath)
	if os.IsNotExist(err) {
		err = s.writeState()
		if err != nil {
			panic(err)
		}
	} else {
		if string(rawState) == "on" {
			s.state = true
		} else {
			s.state = false
		}
	}
}

func (s *kkSwitch) handleRequest(w http.ResponseWriter, req *http.Request) {
	setValue := req.URL.Query().Get("set")

	if len(setValue) > 0 {
		s.handleSet(w, setValue)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"state\":\"%s\"}", makeState(s.state))))
}

func (s *kkSwitch) handleSet(w http.ResponseWriter, value string) {
	var newState bool
	var err error

	switch value {
	case "on":
		newState = true
	case "off":
		newState = false
	case "toggle":
		newState = !s.state
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("I don't know how to do that"))
		return
	}

	if newState != s.state {
		s.state = newState
		err = updateRelay()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = s.writeState()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"ok\":true,\"state\":\"%s\"}", makeState(s.state))))
}

func main() {
	var err error

	mode := &serial.Mode{BaudRate: 9600, Parity: serial.NoParity, StopBits: serial.OneStopBit}
	relay, err = serial.Open(portName, mode)

	if err != nil {
		panic(err)
	}

	err = sendCmd(initCmd)
	if err != nil {
		panic(err)
	}

	lightSwitch.initState()
	motionSwitch.initState()

	err = updateRelay()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/light", lightSwitch.handleRequest)
	http.HandleFunc("/motion", motionSwitch.handleRequest)

	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
