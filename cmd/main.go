package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gitlab.com/verdant-sso/pkg/router"
)

func main() {

	//read env from files instead windows env (applicable for dev only)
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	ver := os.Getenv("BUILD_VER")
	appname := os.Getenv("APP_NAME")
	fmt.Println("Service: ", appname)
	fmt.Println("Version: ", ver)
	server := router.NewRouter()
	server.Run()

}
