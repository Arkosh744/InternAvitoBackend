package main

import (
	"github.com/Arkosh744/InternAvitoBackend/internal/app"
	log "github.com/sirupsen/logrus"

	"os"
)

// @title          ServiceForUserBalanceOperations API
// @version        0.1
// @description    This is Avito Internship Backend Task
// @host           localhost:8080
// @BasePath       /
// @termsOfService http://swagger.io/terms/
// @host           localhost:8080
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
