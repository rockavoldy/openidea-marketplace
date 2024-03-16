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
	if dbUsername == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return ""
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUsername, dbPassword, dbHost, dbPort, dbName)
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
