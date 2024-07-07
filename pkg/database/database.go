package database

import (
	"context"
	"fmt"
	"hng-task-two/pkg/reuseable"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	const maxRetries = 30
	const retryInterval = 2 * time.Second

	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = connect()
		if err == nil {
			return db
		}

		fmt.Printf("Failed to connect to database. Retrying in %v...\n", retryInterval)
		time.Sleep(retryInterval)
	}

	return nil
}

func connect() (*gorm.DB, error) {
	connectionString := fmt.Sprintf("postgres://%v:%v@%v/%v", User, Password, Host, Db_name)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to open database connection:", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Failed to get the *sql.DB object: %s\n", err.Error())
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		fmt.Printf("Failed to ping the database: %s\n", err.Error())
		return nil, err
	}

	fmt.Println("Database connection is active")

	return db, nil
}

var DB *gorm.DB = ConnectDB()

var Host = reuseable.GetEnvVar("DB_HOST")
var Port, _ = strconv.Atoi(reuseable.GetEnvVar("DB_PORT"))
var User = reuseable.GetEnvVar("DB_USER")
var Password = reuseable.GetEnvVar("DB_PASSWORD")
var Db_name = reuseable.GetEnvVar("DB_NAME")
