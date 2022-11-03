package database

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type ConnectionInfo struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func NewPostgresConnection(info ConnectionInfo) (*sql.DB, error) {
	connectString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password)
	log.Println(connectString)
	db, err := sql.Open("postgres", connectString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
