package main

import "http-multiplexer/internal/app"

func main() {
	a := app.NewApp(app.NewConfig())
	a.Run()
}
