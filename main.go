package main

import (
	"context"
	"fmt"
	"net/http"
	"port-logger/app"
	"port-logger/app/config"
	"port-logger/pkg/logger"
)

func main() {
	ctx := context.Background()
	server := app.Server{}

	conf := config.NewConfig()
	conf.Init()

	r := server.Setup()

	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	logger.Infof(ctx, "Listening on: %s", addr)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		logger.Errorf(ctx, "Listening error: %s", err)
	}
}
