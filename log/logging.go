package logging

import (
	"log"
	"os"
)

var Forgelog *log.Logger

func init() {
	Forgelog = log.New(os.Stderr, "", 0)
}
