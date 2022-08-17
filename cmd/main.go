package main

import (
	"go-project/app"
	"time"
)

func main() {
	app := app.Init(
		app.Name("应用名称"),
		app.Version("v1.0"),
	)
	app.Start()
	//time.AfterFunc(time.Second, func() {
	//	app.Stop()
	//})
	time.Sleep(time.Minute)
}
