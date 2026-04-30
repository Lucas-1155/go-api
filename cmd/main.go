package main

import (
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/config"
	"github.com/Comal-Developers1/oracle-grafana-API.git/go-api/router"
)

var (
	logger config.Logger
)

func main() {
	logger = *config.GetLogger("main")
	err := config.Init()

	if err != nil {
		logger.ErrorF("CONFIG INITIALIZATION ERROR: %v", err)
		return
	}

	router.Initialize()
}
