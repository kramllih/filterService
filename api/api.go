package api

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kramllih/filterService/config"
	"github.com/kramllih/filterService/internal/controllers"
	"github.com/kramllih/filterService/internal/database"
	"github.com/kramllih/filterService/internal/logger"
	"github.com/kramllih/filterService/internal/middleware"
)

type HttpConfig struct {
	Host string
	Port int
}

type HttpServer struct {
	Server *gin.Engine
	cfg    HttpConfig
}

func SetupRouter(workpath string, db database.Client, cfg *config.RawConfig, ls string) (Server, error) {

	config := HttpConfig{
		Host: "",
		Port: 80,
	}

	if err := cfg.UnpackRaw(&config); err != nil {
		return nil, err
	}

	ctrl, err := controllers.NewController(ls)
	if err != nil {
		return nil, err
	}

	ctrl.DB = db

	app := gin.New()

	logger := logger.NewLogger("gin")

	app.Use(gin.Recovery())
	app.Use(middleware.RequestID())
	app.Use(middleware.ErrorHandler(logger))
	app.Use(middleware.Logger(logger))

	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":     "RESOURCE_NOT_FOUND",
			"resource": c.Request.RequestURI,
		})
	})

	api := app.Group("/api")
	{
		api.POST("/validate", ctrl.Validate)
		api.GET("/messages", ctrl.AllMessages)
		api.GET("/rejected", ctrl.Rejected)

		ap := api.Group("/approvals")
		{
			ap.POST("/:id/approve", ctrl.Approve)
			ap.POST("/:id/reject", ctrl.Reject)
			ap.GET("", ctrl.AllApprovals)
		}

	}

	h := &HttpServer{
		Server: app,
		cfg:    config,
	}

	return h, nil

}

func (h *HttpServer) Start() error {

	uri := net.JoinHostPort(h.cfg.Host, strconv.Itoa(int(h.cfg.Port)))
	//go func() {
	if err := h.Server.Run(uri); err != nil {
		return fmt.Errorf("Unable to start HTTPS server due to error: %w", err)
	}

	//}()

	return nil
}
