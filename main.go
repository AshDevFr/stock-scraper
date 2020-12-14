package main

import (
	"flag"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"stock_scraper/pkg/scrapers"
	"stock_scraper/pkg/utils"
	"stock_scraper/pkg/websocket"
	"stock_scraper/types"
	"sync"
)

var (
	config       types.Config
	neweggRegex  = regexp.MustCompile(`(?i)newegg\.com`)
	bestbuyRegex = regexp.MustCompile(`(?i)bestbuy\.com`)
	amazonRegex  = regexp.MustCompile(`(?i)amazon\.com`)
)

func init() {
	config = utils.LoadConfig("./config.json")
}

func runScrapers() {
	log.Println("Starting inventory search...")
	// loop for items in config to build and execute http requests

	var waitGroup sync.WaitGroup

	for _, item := range config.Items {
		waitGroup.Add(1)
		go func(item types.Item) {
			defer waitGroup.Done()
			scraperType := ""
			content := ""
			err := ""
			switch {
			case item.Type == "newegg" || neweggRegex.MatchString(item.Url):
				scraperType = "newegg"
				content, err = scrapers.RunNewegg(item)
			case item.Type == "bestbuy" || bestbuyRegex.MatchString(item.Url):
				scraperType = "bestbuy"
				content, err = scrapers.RunBestBuy(item)
			case item.Type == "amazon" || amazonRegex.MatchString(item.Url):
				scraperType = "amazon"
				content, err = scrapers.RunAmazon(item)
			default:
				scraperType = "default"
				content, err = scrapers.RunDefault(item)
			}

			if err != "" {
				websocket.SendUpdateMessage(scraperType, item, "error", err)
			} else {
				websocket.SendUpdateMessage(scraperType, item, "ok", content)
				rules:= config.Rules
				if item.Rules != nil {
					rules = item.Rules
				}
				actions := utils.ApplyRules(rules, "", content)
				for _, action := range actions {
					websocket.SendActionMessage(action, item)
				}
			}
		}(item)
	}
	waitGroup.Wait()
	log.Println("Complete.")
}

func setupCron() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("*/20 * * * * *", func() { runScrapers() })
	c.Start()
}

func setupRouter() (*gin.Engine, *websocket.Hub) {
	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./static", true)))

	// Setup route group for the API
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// Setup websocket
	hub := websocket.GetHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	return router, hub
}

func main() {
	serverEnabled := flag.Bool("s", false, "Enable the server")
	flag.Parse()

	if *serverEnabled {
		setupCron()

		router, _ := setupRouter()

		router.Run(":5000")
	} else {
		runScrapers()
	}
}
