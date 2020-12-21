package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/webview/webview"
	"net/http"
	configUtils "stock_scraper/internal/config"
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

	configUtils.LoadUserAgents("./user_agents.txt")

	maxConcurrency := defaultMaxConcurrency
	if config.Concurrency != nil {
		maxConcurrency = *config.Concurrency
	}
	parallelChan = make(chan bool, maxConcurrency)
}

func loadWebview(addr string) {
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Stock scraper")
	w.SetSize(1280, 1024, webview.HintNone)
	w.Navigate(addr)
	w.Run()
}

func main() {
	serverOpt := flag.Bool("s", false, "Enable the server")
	uiOpt := flag.Bool("ui", false, "Enable the ui")
	verboseOpt := flag.Bool("v", false, "Verbose")
	watchOpt := flag.Bool("w", false, "Watch for changes (Not required if the server is enabled)")
	configFile := flag.String("f", "config.json", "Config file")

	flag.Parse()

	config = configUtils.LoadConfig(*configFile)

	if *verboseOpt {
		log.SetLevel(log.DebugLevel)
		log.Info("Mode verbose enabled")
	}

	if *serverOpt {
		log.Info("Web server enabled")
		setupCron()

		router, _ := setupRouter()

		addr := "localhost:5000"
		if *uiOpt {
			go router.Run(addr)
			loadWebview(fmt.Sprintf("http://%s", addr))
		} else {
			router.Run(addr)
		}
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
