package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"net/http"
	config2 "stock_scraper/internal/config"
	"stock_scraper/types"
	"time"
)

var (
	defaultMaxConcurrency = 25
	parallelChan          chan bool
	config                types.Config
)

func init() {
	http.DefaultClient.Timeout = time.Second * 300

	config = config2.LoadConfig("./config.json")
	config2.LoadUserAgents("./user_agents.txt")

	maxConcurrency := defaultMaxConcurrency
	if config.Concurrency != nil {
		maxConcurrency = *config.Concurrency
	}
	parallelChan = make(chan bool, maxConcurrency)
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
