package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var (
	logger         *slog.Logger
	oneLogger      sync.Once
	dbPool         *sqlx.DB
	oneDBPool      sync.Once
	mySQLDBPool    *sqlx.DB
	oneMySQLDBPool sync.Once
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetMySQLDB() *sqlx.DB {
	oneMySQLDBPool.Do(
		func() {
			var err error
			dbUrl := os.Getenv("MYSQL_DATABASE_URL")
			connStr := fmt.Sprintf("%v", dbUrl)
			mySQLDBPool, err = sqlx.Open("mysql", connStr)
			if err != nil {
				panic(err)
			}
			mySQLDBPool.SetConnMaxLifetime(0)
			mySQLDBPool.SetMaxIdleConns(3)
			mySQLDBPool.SetMaxOpenConns(3)
		},
	)
	return mySQLDBPool
}

func GetDB() *sqlx.DB {
	oneDBPool.Do(
		func() {
			var err error
			dbUrl := os.Getenv("DATABASE_URL")
			connStr := fmt.Sprintf("%v", dbUrl)
			dbPool, err = sqlx.Open("postgres", connStr)
			if err != nil {
				panic(err)
			}
			dbPool.SetConnMaxLifetime(0)
			dbPool.SetMaxIdleConns(3)
			dbPool.SetMaxOpenConns(3)
		},
	)
	return dbPool
}

func GetLogger() *slog.Logger {
	oneLogger.Do(
		func() {
			logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		},
	)
	return logger
}
