package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lorsanstand/gontrol/internal/hub/app"
)

func main() {
	app := app.NewApp()

	app.Run()
}
