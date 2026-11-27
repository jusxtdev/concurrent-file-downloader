package main

import (
	"log"

	"github.com/jusxtdev/concurrent-file-download/cmd"
	"github.com/jusxtdev/concurrent-file-download/internal/controller"
)

func main() {
	controller := controller.NewController()
	cli := cmd.NewCLIApplication(controller)
	cli.Setup()
	err := cli.Loop()
	if err != nil {
		log.Fatal(err)
	}
	// TODO :- add a cleanup function which deletes empty folders in `temp/`
}
