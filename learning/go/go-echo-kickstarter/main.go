package main

import (
	"go-echo-kickstarter/internal/app"
	"time"
)

func main() {
	a := app.New()
	defer a.Shutdown(5 * time.Second)

	a.Run()
}
