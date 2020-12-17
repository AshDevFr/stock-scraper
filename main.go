package main

import (
	"flag"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"stock_scraper/pkg/browser"
	config2 "stock_scraper/pkg/config"
	"stock_scraper/pkg/process"
	"stock_scraper/pkg/state"
	"stock_scraper/pkg/websocket"
	"stock_scraper/types"
	"sync"
	"time"
)

var (
	config types.Config
)

func init() {
	http.DefaultClient.Timeout = time.Second * 120
	config = config2.LoadConfig("./config.json")
	config2.LoadUserAgents("./user_agents.txt")
}

func runScraper(item types.Item) {
	scraperType := item.Parser.Label()
	content, err := item.Parser.Run(item)

	if err != "" {
		websocket.SendUpdateMessage(scraperType, item, "error", err)
	} else {
		websocket.SendUpdateMessage(scraperType, item, "ok", content)
		rules := config.DefaultConfig.Rules
		if item.Config.Rules != nil {
			rules = item.Config.Rules
		}
		previousContent := state.GetContent(item.Uuid)
		actions := process.ApplyRules(rules, previousContent, content)
		state.SetContent(item.Uuid, content)

		if len(actions) > 0 {
			if *item.Config.OpenLinks && item.Url != "" {
				state.ShouldRunAlert(item.Uuid, func() {
					if item.AddToCartUrl != "" {
						browser.Open(item.AddToCartUrl)
					} else {
						browser.Open(item.Url)
					}
				})
				log.Info("Opening link")
			}
			for _, action := range actions {
				websocket.SendActionMessage(action, item)
			}
		}
	}
}

func runAllScrapers() {
	log.Debug("Running the scrapers: start")
	// loop for items in config to build and execute http requests

	var waitGroup sync.WaitGroup

	for _, item := range config.Items {
		waitGroup.Add(1)
		go func(item types.Item) {
			defer waitGroup.Done()
			runScraper(item)
		}(item)
	}
	waitGroup.Wait()
	log.Debug("Running the scrapers: done")
}

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
	serverOpt := flag.Bool("s", false, "Enable the server")
	verboseOpt := flag.Bool("v", false, "Verbose")
	watchOpt := flag.Bool("w", false, "Watch for changes (Not required if the server is enabled)")

	flag.Parse()

	if *verboseOpt {
		log.SetLevel(log.DebugLevel)
		log.Info("Mode verbose enabled")
	}

	if *serverOpt {
		log.Info("Web server enabled")
		setupCron()

		router, _ := setupRouter()

		router.Run(":5000")
	} else {
		if *watchOpt {
			setupCron()

			for {
				time.Sleep(10 * time.Second)
			}
		} else {
			runAllScrapers()
		}
	}
}
