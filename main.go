package main

import (
	"github.com/emvi/logbuch"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/web/handlers"
	"github.com/muety/mailwhale/web/routes/api"
	"github.com/rs/cors"
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
	recoverMiddleware := ghandlers.RecoveryHandler()
	loggingMiddleware := handlers.NewLoggingMiddleware(logbuch.Info, []string{})

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: config.Web.CorsOrigins,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Configure routing
	router := mux.NewRouter().StrictSlash(true)
	router.Use(recoverMiddleware, loggingMiddleware)

	// Handlers
	api.NewHealthHandler().Register(router)
	api.NewMailHandler().Register(router)
	api.NewClientHandler().Register(router)
	api.NewTemplateHandler().Register(router)

	handler := corsHandler.Handler(router)

	// Static routes
	router.PathPrefix("/").Handler(handlers.SPAHandler{
		StaticPath: "./webui/public",
		IndexPath:  "index.html",
	})

	listen(handler, config)
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

		user, err := userService.Create(&types.Signup{
			Email:    u.Email,
			Password: u.Password,
		})
		if err != nil {
			logbuch.Fatal("failed to create seed user '%s': %v", u.Email, err)
		}
		logbuch.Info("created seed user '%s'", user.ID)
	}
}
