package app

import (
	"context"
	"fmt"
	"github.com/Arkosh744/InternAvitoBackend/internal/config"
	"github.com/Arkosh744/InternAvitoBackend/internal/repository"
	"github.com/Arkosh744/InternAvitoBackend/internal/service"
	"github.com/Arkosh744/InternAvitoBackend/internal/transport/rest"
	"github.com/Arkosh744/InternAvitoBackend/pkg/database"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title ServiceForUserBalanceOperations API
// @version 0.1
// @description This is Avito Internship Backend Task

// @host localhost:8080
// @BasePath /

// @termsOfService http://swagger.io/terms/
// @host localhost:8080
func Run() error {
	cfg, err := config.New("configs")
	if err != nil {
		return err
	}

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Username: cfg.DBUser,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
		Password: cfg.DBPassword,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	usersRepo := repository.NewUsers(db)
	userService := service.NewUsersService(usersRepo)
	handler := rest.NewHandler(userService)

	srv := &http.Server{
		Addr:         ":" + cfg.SrvPort,
		Handler:      handler.InitRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Info("Starting Server at port " + cfg.SrvPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	// GRACEFUL SHUTDOWN with 5 seconds timeout BELOW
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("timeout of 1 seconds.")
	<-ctx.Done()
	log.Println("Server exiting")
	return nil
}
