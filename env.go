package main

import "os"

var isDev = os.Getenv("APP_ENV") == "development"

const devTimeout = 500
