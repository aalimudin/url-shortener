package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	redisRepository := newRedisRepository("localhost:6379", "", 0)
	shortenService := newShortenService(redisRepository)

	handler := newHandler(shortenService)
	e := echo.New()

	e.POST("/shorten", handler.shortenHandler)
	e.GET("/:slug", handler.getShortenURLHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
