package main

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"stock_scraper/internal/websocket"
)

func setupCron() {
	log.Debug("Setting up the cron jobs: start")
	c := cron.New(cron.WithSeconds())

	for _, item := range config.Items {
		it := item
		c.AddFunc(item.Config.Cron, func() {
			runScraper(it)
		})
	}
	c.Start()
	log.Debug("Setting up the cron jobs: done")
}

func setupRouter() (*gin.Engine, *websocket.Hub) {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": config,
			})
		})
	}

	// Setup websocket
	hub := websocket.GetHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	router.Use(gin.WrapH(http.FileServer(pkger.Dir("/static"))))

	return router, hub
}
