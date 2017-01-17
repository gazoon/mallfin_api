package main

import (
	"github.com/gazoon/httprouter"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte(ps.ByName("id")))
}
func main() {
	r := httprouter.New()
	r.GET("/shops/:id/", handler)
	panic(http.ListenAndServe(":8080", r))
}
