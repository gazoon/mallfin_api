package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"mallfin_api/config"
	"mallfin_api/redisdb"
	"mallfin_api/utils"
	"net/http"
)

type User struct {
	Age  int `form:"user_age"`
	Name string
}

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u := User{}
	log.Infof("User: %+v", u)
	m := utils.NewDistributedMutex("kaaa")
	m.Lock()
}

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	var configPath string
	flag.StringVar(&configPath, "conf", "", "Path to json config file.")
	flag.Parse()
	if configPath == "" {
		log.WithField("location", "main").Panic("Cannot start without path to config")
	}
	config.Initialization(configPath)
	redisdb.Initialization()
	defer redisdb.Close()
	r := httprouter.New()
	r.GET("/", handler)
	n := negroni.New()

	n.Use(&negroni.Recovery{
		Logger: log.StandardLogger(),
	})
	n.Use(&negroni.Logger{ALogger: log.StandardLogger()})
	n.UseHandler(r)
	http.ListenAndServe(":8080", n)
}
