package main

import (
	"github.com/emvi/logbuch"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/web/middleware"
	"github.com/muety/mailwhale/web/routes/api"
	"github.com/timshannon/bolthold"
	"net/http"
	"time"
)

var (
	config      *conf.Config
	store       *bolthold.Store
	userService *service.UserService
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
	userService = service.NewUserService()

	initDefaults()

	// Global middlewares
	recoverMiddleware := handlers.RecoveryHandler()
	loggingMiddleware := middleware.NewLoggingMiddleware(logbuch.Info, []string{})

	// Configure routing
	router := mux.NewRouter()
	router.Use(recoverMiddleware, loggingMiddleware)

	// Handlers
	api.NewHealthHandler().Register(router)
	api.NewMailHandler().Register(router)
	api.NewClientHandler().Register(router)

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
		logbuch.Info("web server started, listening on %s", config.Web.ListenV4)
		if err := s4.ListenAndServe(); err != nil {
			logbuch.Fatal("failed to start web server: %v", err)
		}
	}()

	<-make(chan interface{}, 1)
}

func initDefaults() {
	for _, u := range config.Security.SeedUsers {
		if user, err := userService.GetById(u.Email); err == nil {
			user.Password = u.Password
			user, err = userService.Update(user)
			if err != nil {
				logbuch.Fatal("failed to update user '%s': %v", u.Email, err)
			}
			logbuch.Info("updated user '%s'", user.ID)
			continue
		}

		user, err := userService.Create(&u)
		if err != nil {
			logbuch.Fatal("failed to create seed user '%s': %v", u.Email, err)
		}
		logbuch.Info("created seed user '%s'", user.ID)
	}
}
