package main

import (
	"github.com/emvi/logbuch"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/web/routes/api"
	"net/http"
	"time"
)

var (
	config *conf.Config
)

func main() {
	config = conf.Load()

	// Set log level
	if config.IsDev() {
		logbuch.SetLevel(logbuch.LevelDebug)
	} else {
		logbuch.SetLevel(logbuch.LevelInfo)
	}

	// Services
	sendService := service.NewSendService()

	// Configure routing
	router := http.NewServeMux()

	// Handlers
	api.NewHealthHandler().Register(router)
	api.NewMailHandler(sendService).Register(router)

	listen(router, config)
}

func listen(handler http.Handler, config *conf.Config) {
	var s4 *http.Server

	s4 = &http.Server{
		Handler:      handler,
		Addr:         config.Web.ListenV4,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logbuch.Info("Web server started. Listening on %s", config.Web.ListenV4)
		if err := s4.ListenAndServe(); err != nil {
			logbuch.Fatal("failed to start web server: %v", err)
		}
	}()

	<-make(chan interface{}, 1)
}
