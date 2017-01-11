package main

import (
	"fmt"
	"mallfin_api/config"
	"mallfin_api/db"
	"mallfin_api/handlers"
	"mallfin_api/redisdb"
	"net/http"
	"runtime/debug"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/httprouter"
	"github.com/urfave/negroni"
)

func recoveryMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal server error")
			log.WithField("location", "recovery middleware").Warnf("Panic recovered: %s", err)
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
	r.GET("/shops/", handlers.ShopsList)
	r.GET("/shops/:id/", handlers.ShopDetails)

	n := negroni.New()
	n.Use(&negroni.Logger{ALogger: log.StandardLogger()})
	n.UseFunc(recoveryMiddleware)
	n.UseHandler(r)
	mainLogger.Infof("Starting server on port %d", config.Port())
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port()), n)
	if err != nil {
		mainLogger.Panicf("Cannot run server: %s", err)
	}
}
