package controllers

import (
	"github.com/kramllih/filterService/internal/database"
	"github.com/kramllih/filterService/internal/httpClient"
	"github.com/kramllih/filterService/internal/logger"
)

type Config struct {
	Host string
}

type Controller struct {
	httpClient *httpClient.HTTP
	DB         database.Client
	log        *logger.Logger
}

func NewController(ls string) (*Controller, error) {

	config := Config{
		Host: "http://localhost:8081",
	}

	if ls != "" {
		config.Host = ls
	}

	http := httpClient.NewHTTP()
	http.SetURI(config.Host)

	return &Controller{
		log:        logger.NewLogger("controller"),
		httpClient: http,
	}, nil
}
