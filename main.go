package main

import (
	"github.com/namsral/flag"

	"github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/thatarchguy/webconsul/controllers"
)

var config struct {
	port       string
	consulAddr string
	consulDC   string
}

func main() {
	loadConfig()

	kvConfig := api.DefaultConfig()
	kvConfig.Address = config.consulAddr
	kvConfig.Datacenter = config.consulDC

	client, err := api.NewClient(kvConfig)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	controller := controllers.New(client.KV())

	e.GET("/health", func(ctx echo.Context) error {
		return ctx.JSON(200, map[string]string{"status": "ok"})
	})

	v1 := e.Group("v1")
	v1.GET("/consul", controller.Consul)

	e.Logger.Fatal(e.Start(":" + config.port))
}

func loadConfig() {

	flag.StringVar(&config.consulAddr, "addr", "https://localhost", "Consul HTTP Addr")
	flag.StringVar(&config.consulDC, "dc", "", "consul datacenter, uses local if blank")
	flag.StringVar(&config.port, "port", "8080", "port to run webserver on")
	flag.Parse()
}
