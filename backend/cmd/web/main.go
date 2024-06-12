package main

import (
	_ "CoolUrlShortener/docs"
	"CoolUrlShortener/internal/app"
)

//	@title			CoolURLShortener API
//	@version		1.0
//	@description	API Server for shorten urls

// @BasePath	/
func main() {
	app.Run()
}
