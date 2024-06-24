package main

import (
	"github.com/shashank-sharma/backend/internal/app"
	_ "github.com/shashank-sharma/backend/migrations"
)

func main() {
	application := app.New()
	application.Start()
}
