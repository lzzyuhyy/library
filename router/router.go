package router

import (
	"github.com/gin-gonic/gin"
	"library/api"
)

func Router(g *gin.Engine) {
	router := g.Group("api")

	router.POST("book/list", api.BookList)
	router.POST("book/borrow", api.BorrowBook)
	router.POST("pay", api.Pay)
	router.POST("notify", api.Notify)
}
