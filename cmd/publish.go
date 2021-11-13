package main

import (
	"github.com/reliqarts/go-docweaver"
	"os"
)

func main() {
	args := os.Args[1:]
	logger := docweaver.GetLoggerSet()
	shouldUpdate := true

	if len(args) < 2 {
		logger.Err.Fatal("One or more arguments missing. Usage: `publish productName productSource [shouldUpdate=true]`")
	}

	// if third arg is passed, and is false, disable updates
	if len(args) == 3 && (args[2] == "false" || args[2] == "f") {
		shouldUpdate = false
	}

	publisher := docweaver.GetPublisher()
	publisher.Publish(args[0], args[1], shouldUpdate)
}
