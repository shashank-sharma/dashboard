package main

import (
	"github.com/shashank-sharma/backend/internal/app"
)

func main() {
	application := app.New()
	application.Start()
}
