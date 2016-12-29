package main

import (
	"fmt"

	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"

	"net/http"
)

type User struct {
	Age  int `form:"user_age"`
	Name string
}

func handler(ctx *iris.Context) {
	u := User{}
	fmt.Println("hello worlddddd")
	err := ctx.ReadForm(&u)
	if err != nil {
		ctx.Log("ReadForm err: %s", err)
		panic(err)
	}
	ctx.Log("User: %+v", u)
	ctx.JSON(http.StatusOK, map[string]interface{}{"foo": "bar"})
}
func main() {
	iris.Use(recovery.New())
	iris.OnError(iris.StatusInternalServerError, func(ctx *iris.Context) {
		ctx.JSON(500, map[string]interface{}{"er": 2})
	})
	iris.Use(logger.New())
	iris.Get("/", handler)
	iris.Listen(":8080")
}
