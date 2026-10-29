package main

import (
	"log"
	"os"

	"github.com/lorsanstand/gontrol/internal/agent"
)

func main() {
	app, err := agent.NewApp()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(app.Run())
}
