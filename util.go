package restclient

import (
	"log"
	"os"
)

var defaultLogger = log.New(os.Stdout, "", 0)
