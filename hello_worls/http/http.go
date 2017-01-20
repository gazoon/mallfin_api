package main

import (
	"net/http"

	"github.com/gazoon/httprouter"
	"gopkg.in/redis.v5"
)

var client *redis.Client

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s, err := client.Get("some_common_key").Result()
	if err != nil {
		panic(err)
	}
	w.Write([]byte(s))
}
func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 10,
		DB:       2,
	})
	err := client.Ping().Err()
	if err != nil {
		panic(err)
	}
	r := httprouter.New()
	r.GET("/shops/:id/", handler)
	panic(http.ListenAndServe(":8080", r))
}
