package main

import (
	"github.com/condrowiyono/ruangtengah-api/app"
	"github.com/condrowiyono/ruangtengah-api/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":9000")
}
