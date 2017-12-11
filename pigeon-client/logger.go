package main

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)

const (
	DEBUG = false
)

func init() {
	mode := log.Ldate | log.Ltime
	if DEBUG {
		mode = log.Ldate | log.Ltime | log.Lshortfile
	}
	logger = log.New(os.Stdout, "", mode)

	logger.SetPrefix("[pigeon client] ")
}
