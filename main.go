package main

import (
	"github.com/condrowiyono/living-rooms-api/app"
	"github.com/condrowiyono/living-rooms-api/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":9000")
}
