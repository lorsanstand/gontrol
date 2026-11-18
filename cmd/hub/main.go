package main

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lorsanstand/gontrol/internal/hub/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(app.Run())
}
