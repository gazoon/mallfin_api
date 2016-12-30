package main

import (
	"flag"
	"fmt"
	"mallfin_api/config"
	"mallfin_api/db"
	"mallfin_api/redisdb"
	"net/http"
	"runtime/debug"

	log "github.com/Sirupsen/logrus"
	"github.com/gazoon/httprouter"
	"github.com/urfave/negroni"
)

type User struct {
	Age  int `form:"user_age"`
	Name string
}

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	valInt, err := ps.ByNameInt("foo")
	if err != nil {
		log.Warn("NOT INT")
		return
	}
	log.Info(valInt == 3)
	log.Info(valInt)
	u := User{}
	log.Infof("User: %+v", u)

}
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

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	mainLogger := log.WithField("location", "main")
	var configPath string
	flag.StringVar(&configPath, "conf", "", "Path to json config file.")
	flag.Parse()

	if configPath == "" {
		mainLogger.Panic("Cannot start without path to config")
	}
	config.Initialization(configPath)

	redisdb.Initialization()
	defer redisdb.Close()

	db.Initialization()
	defer db.Close()

	r := httprouter.New()
	r.GET("/:foo/", handler)

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
