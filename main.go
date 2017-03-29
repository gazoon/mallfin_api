package main

import (
	"fmt"
	"mallfin_api/config"
	"mallfin_api/db"
	"mallfin_api/handlers"
	"mallfin_api/redisdb"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"

	"flag"

	"mallfin_api/logging"
	"mallfin_api/tracing"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/httprouter"
	"github.com/urfave/negroni"
)

var logger = logging.WithPackage("main")

func recoveryMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal server error")
			if _, isLogPanic := err.(*log.Entry); !isLogPanic {
				log.WithField("location", "recovery middleware").Errorf("Panic recovered: %s", err)
			}
			if config.Debug() {
				debug.PrintStack()
			}
		}
	}()
	next(w, r)
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "conf", "", "Path to json config file.")
	flag.Parse()

	config.Initialization(configPath)
	logging.Initialization()

	redisdb.Initialization()
	defer redisdb.Close()

	db.Initialization()
	defer db.Close()

	r := httprouter.New()
	r.GET("/malls/", handlers.MallsList)
	r.GET("/malls/:id/", handlers.MallDetails)
	r.GET("/current_mall/", handlers.CurrentMall)
	r.GET("/current_city/", handlers.CurrentCity)
	r.GET("/shops_in_malls/", handlers.ShopsInMalls)
	r.GET("/search/", handlers.Search)
	r.GET("/shops/", handlers.ShopsList)
	r.GET("/shops/:id/", handlers.ShopDetails)
	r.GET("/categories/", handlers.CategoriesList)
	r.GET("/categories/:id/", handlers.CategoryDetails)
	r.GET("/cities/", handlers.CitiesList)

	n := negroni.New()
	n.UseFunc(recoveryMiddleware)
	//c := cors.New(cors.Options{AllowedOrigins: []string{"*"}})
	//n.Use(c)
	n.UseFunc(tracing.Middleware)
	n.UseFunc(logging.Middleware)
	n.UseHandler(r)
	if config.Debug() {
		go func() {
			err := http.ListenAndServe(":6060", nil)
			if err != nil {
				logger.Panicf("Cannot run profiler server: %s", err)
			}
		}()
	}
	logger.Infof("Starting server on port %d", config.Port())
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port()), n)
	if err != nil {
		logger.Panicf("Cannot run server: %s", err)
	}
}
