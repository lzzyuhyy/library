package main

import (
	"github.com/gin-gonic/gin"
	"library/initial"
	"library/router"
)

func main() {
	err := initial.Initial()
	if err != nil {
		panic(err)
	}

	g := gin.Default()
	g.Use()
	router.Router(g)

	err = g.Run(":9999")
	if err != nil {
		return
	}
}
