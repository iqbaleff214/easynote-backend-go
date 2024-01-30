package main

import "os"

type config struct {
	mysqlUri  string
	jwtSecret string
	version   string
}

var appConfig config

func initConfig() {
	mysqlUri := os.Getenv("MYSQL_URI")
	if mysqlUri == "" {
		mysqlUri = "root:@tcp(127.0.0.1:3306)/easynote?parseTime=true"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "easynotejwtsecret123"
	}
	
	version := os.Getenv("VERSION")
	if version == "" {
		version = "1"
	}

	appConfig = config{mysqlUri, jwtSecret, version}
}
