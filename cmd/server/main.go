package main

import (
	"go-printer/internal/app"
	"log"
)

func main() {
	log.Println("Hello world")

	a := app.NewApp()
	a.Run()
}
