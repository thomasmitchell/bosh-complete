package main

import (
	"fmt"
	"os"
)

var log logger

type logger struct {
	on bool
	f  *os.File
}

func (l *logger) TurnOn() {
	var err error

	if _, err := os.Stat("/tmp/bosh_complete"); err != nil {
		if !os.IsNotExist(err) {
			panic("Could not stat log dir: " + err.Error())
		}
		err = os.Mkdir("/tmp/bosh_complete", 0775)
		if err != nil {
			panic("Could not make log dir: " + err.Error())
		}
	}

	l.f, err = os.OpenFile("/tmp/bosh_complete/log.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic("Could not open logging file" + err.Error())
	}

	l.on = true
}

func (l logger) Write(f string, args ...interface{}) {
	if !l.on {
		return
	}
	l.f.Write([]byte(fmt.Sprintf("%s\n", fmt.Sprintf(f, args...))))
}
