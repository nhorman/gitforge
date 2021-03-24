package logging

import (
	"log"
	"os"
	"path/filepath"
)

var Forgelog *log.Logger

func init() {
	Forgelog = log.New(os.Stderr, "", 0)
}

func SuppressLog() {
	dnf, _ := os.Open(os.DevNull)
	Forgelog = log.New(dnf, "", 0)
}

func LogToFile(logfile string) {
	abspath, _ := filepath.Abs(logfile)
	fd, _ := os.Open(abspath)
	Forgelog = log.New(fd, "", 0)
}
