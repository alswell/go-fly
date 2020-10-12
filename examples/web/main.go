package main

import (
	"github.com/alswell/go-fly/web"

	_ "github.com/alswell/go-fly/examples/web/api"
)

func init() {
	web.RegisterRouter(web.GET, "/",
		func() string {
			return "welcome"
		},
	)
}

func main() {
	web.StartWebServer(8080)
}
