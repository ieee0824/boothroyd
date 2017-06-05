package main

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"github.com/jobtalk/hawkeye/controllers/queue"
)

func main() {
	app := iris.New()
	app.Adapt(
		iris.DevLogger(),
		httprouter.New(),
		cors.New(cors.Options{AllowedOrigins: []string{"*"}}),
	)

	app.Post("/queue/enqueue", queue.Enqueue)
	app.Post("/queue/check", queue.Check)

	go queue.Run()

	app.Listen(":8080")
}

