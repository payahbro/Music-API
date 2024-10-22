package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"music-echo/api/handler"
	"net/http"
)

func Init(e *echo.Echo, tracksHandler handler.TracksHandler, userHandler handler.UserHandler) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/v1/healthcheck", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadGateway)
	})

	// tracks
	e.GET("/v1/tracks", tracksHandler.GetAllTracks)
	e.POST("/v1/tracks", tracksHandler.CreateTracks)
	e.GET("/v1/tracks/:tracksId", tracksHandler.GetTracksByID)
	e.PATCH("/v1/tracks/:tracksId", tracksHandler.UpdateTracks)
	e.DELETE("/v1/tracks/:tracksId", tracksHandler.DeleteTracks)
	e.PATCH("/v1/tracks/:tracksId/like", tracksHandler.LikeTracks)

	// users
	e.POST("/v1/users", userHandler.CreateUser)
	e.PUT("v1/users/activated", userHandler.ActivateUser)
}
