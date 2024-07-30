package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type handler struct {
	shortenService shortenServiceInterface
}

type shortenRequest struct {
	Url        string `json:"url"`
	CustomSlug string `json:"customSlug,omitempty"`
}

type shortenResponse struct {
	Url string `json:"shortenURL"`
}

func newHandler(shortenService shortenServiceInterface) *handler {
	return &handler{
		shortenService: shortenService,
	}
}

func (h *handler) shortenHandler(c echo.Context) error {
	req := shortenRequest{}
	ctx := context.Background()

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	shortenURL, err := h.shortenService.ShortenURL(ctx, req.Url, req.CustomSlug)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	res := shortenResponse{
		Url: shortenURL,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) getShortenURLHandler(c echo.Context) error {
	slug := c.Param("slug")
	ctx := context.Background()

	actualURL, err := h.shortenService.GetShortenURL(ctx, slug)
	if err != nil {
		log.Error(err.Error())
		return c.HTML(http.StatusNotFound, err.Error())
	}

	return c.Redirect(http.StatusPermanentRedirect, actualURL)
}
