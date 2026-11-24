package main

import (
	"log"

	"github.com/jusxtdev/concurrent-file-download/cmd"
)

func main(){
	cli := cmd.NewCLIApplication()
	cli.Setup()
	err := cli.Loop()
	if err != nil {
		log.Fatal(err)
	}
}