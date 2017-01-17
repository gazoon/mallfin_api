package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func handler(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte(ctx.UserValue("id").(string)))
}
func main() {
	r := fasthttprouter.New()
	r.GET("/shops/:id/", handler)
	panic(fasthttp.ListenAndServe(":8080", r.Handler))
}
