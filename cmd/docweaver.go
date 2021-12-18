package main

import (
	"github.com/reliqarts/go-docweaver"
	"log"
	"os"
)

var publisher = docweaver.GetPublisher()

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("One or more arguments missing. Usage: `docweaver (publish productName productSource [shouldUpdate=true])|(update [...productNames])`")
	}

	action := args[0]
	switch action {
	case "update":
		update(args[1:]...)
		return
	case "publish":
		publish(args[1:]...)
		return
	default:
		log.Fatalf("Invalid action given: `%s`. Must be 'publish' or 'update'.", action)
	}
}

func update(args ...string) {
	if len(args) == 0 {
		log.Println("No product names given for update. All products will be updated.")
		publisher.UpdateAll()
		return
	}

	publisher.Update(args...)
}

func publish(args ...string) {
	shouldUpdate := true

	if len(args) < 2 {
		log.Fatal("One or more arguments missing for publish action. Usage: `publish productName productSource [shouldUpdate=true]`")
	}

	// if third arg is passed, and is false, disable updates
	if len(args) == 3 && (args[2] == "false" || args[2] == "f") {
		shouldUpdate = false
	}

	publisher.Publish(args[0], args[1], shouldUpdate)
}
