package logging

import (
	"log"
	"os"
)

var Forgelog *log.Logger

func init() {
	Forgelog = log.New(os.Stderr, "", 0)
}

func SuppressLog() {
	dnf, _ := os.Open(os.DevNull)
	Forgelog = log.New(dnf, "", 0)
}
