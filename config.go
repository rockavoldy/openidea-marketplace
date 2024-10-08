package main

import (
	"fmt"
	"os"
)

func DbConnStr() string {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	var dbParams string
	if dbUsername == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return ""
	}
	if os.Getenv("ENV") == "production" {
		dbParams = "sslmode=verify-full&sslrootcert=ap-southeast-1-bundle.pem"
	} else {
		dbParams = "sslmode=disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dbUsername, dbPassword, dbHost, dbPort, dbName, dbParams)
}

func JWTSecretKey() string {
	jwtSeretKey := os.Getenv("JWT_SECRET")
	if jwtSeretKey == "" {
		return "justrandomstringfordefaultpurposes"
	}

	return jwtSeretKey
}

var (
	S3_ID          string
	S3_SECRET_KEY  string
	S3_BUCKET_NAME string
	S3_REGION      string
)

func LoadS3FromEnv() {
	S3_ID = os.Getenv("S3_ID")
	S3_SECRET_KEY = os.Getenv("S3_SECRET_KEY")
	S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME")
	S3_REGION = os.Getenv("S3_REGION")
}
