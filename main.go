package main

import (
	"github.com/emvi/logbuch"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/web/middleware"
	"github.com/muety/mailwhale/web/routes/api"
	"github.com/timshannon/bolthold"
	"net/http"
	"time"
)

var (
	config        *conf.Config
	store         *bolthold.Store
	sendService   *service.SendService
	clientService *service.ClientService
)

func main() {
	config = conf.Load()
	store = conf.LoadStore(config.Store.Path)
	defer store.Close()

	// Set log level
	if config.IsDev() {
		logbuch.SetLevel(logbuch.LevelDebug)
	} else {
		logbuch.SetLevel(logbuch.LevelInfo)
	}

	// Services
	sendService = service.NewSendService()
	clientService = service.NewClientService()

	initDefaults()

	// Global middlewares
	recoverMiddleware := middleware.NewRecoverMiddleware()
	loggingMiddleware := middleware.NewLoggingMiddleware(logbuch.Info, []string{})
	baseChain := alice.New(recoverMiddleware, loggingMiddleware)

	// Configure routing
	router := httprouter.New()

	// Handlers
	api.NewHealthHandler().Register(router, &baseChain)
	api.NewMailHandler(sendService, clientService).Register(router, &baseChain)
	api.NewClientHandler(clientService).Register(router, &baseChain)

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

func initDefaults() {
	clients, err := clientService.GetAll()
	if err != nil {
		logbuch.Fatal("failed to fetch clients initially: %v", err)
	}

	if len(clients) == 0 {
		client, err := clientService.Create(&types.Client{
			Name:        conf.DefaultClientName,
			Permissions: []string{conf.PermissionSendMail, conf.PermissionManageClient},
		})
		if err != nil {
			logbuch.Fatal("failed to initialize default client: %v", err)
		}
		logbuch.Info("Created default client with name '%s'", conf.DefaultClientName)
		logbuch.Info(">>> API key: %s <<<", client.ApiKey)
	}
}
