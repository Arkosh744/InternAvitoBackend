package main

import (
	"github.com/Arkosh744/InternAvitoBackend/internal/app"
	log "github.com/sirupsen/logrus"

	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
