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

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

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
func someTests() {
}
func main() {
	someTests()
	mainLogger := log.WithField("location", "main")

	config.Initialization()

	redisdb.Initialization()
	defer redisdb.Close()

	db.Initialization()
	defer db.Close()

	r := httprouter.New()
	r.GET("/malls/", handlers.MallsList)
	r.GET("/malls/:id/", handlers.MallDetails)
	r.GET("/current_mall/", handlers.CurrentMall)
	r.GET("/shops_in_malls/", handlers.ShopsInMalls)
	r.GET("/search/", handlers.Search)
	r.GET("/shops/", handlers.ShopsList)
	r.GET("/shops/:id/", handlers.ShopDetails)
	r.GET("/categories/", handlers.CategoriesList)
	r.GET("/categories/:id/", handlers.CategoryDetails)
	r.GET("/cities/", handlers.CitiesList)

	n := negroni.New()
	c := cors.New(cors.Options{AllowedOrigins: []string{"*"}})
	n.Use(c)
	n.UseFunc(recoveryMiddleware)
	if config.AccessLog() {
		l := negroni.NewLogger()
		l.ALogger = log.StandardLogger()
		n.Use(l)
	}
	n.UseHandler(r)
	if config.Debug() {
		go func() {
			err := http.ListenAndServe(":6060", nil)
			if err != nil {
				mainLogger.Panicf("Cannot run profiler server: %s", err)
			}
		}()
	}
	mainLogger.Infof("Starting server on port %d", config.Port())
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port()), n)
	if err != nil {
		mainLogger.Panicf("Cannot run server: %s", err)
	}
}
