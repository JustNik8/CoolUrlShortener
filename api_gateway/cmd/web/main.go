package main

import (
	_ "api_gateway/docs"
	"api_gateway/internal/app"
)

//	@title			CoolURLShortener API
//	@version		1.0
//	@description	API Server for shorten urls

func main() {
	app.Run()
}
