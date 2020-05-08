package main

import (
	"fmt"
	"lottery/bootstrap"
	"lottery/conf"
	"lottery/web/middleware/identity"
	"lottery/web/routes"
)

var port = 8080

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("Go抽奖系统", "bairn")
	app.Bootstrap()
	app.Configure(identity.Configure, routes.Configure)
	return app
}

func main() {
	if port == 8080 {
		conf.RunningCrontabService = true
	}

	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}
