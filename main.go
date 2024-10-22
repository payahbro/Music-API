package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"music-echo/api/handler"
	"music-echo/api/repository"
	"music-echo/api/router"
	"music-echo/utils"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// CONFIG
	// Echo
	e := echo.New()
	// Database
	Db := utils.OpenDB()
	defer Db.Close()

	// SECONDARY
	// Validator
	validators := validator.New()
	// Mailer
	mailer := utils.NewMailer("sandbox.smtp.mailtrap.io", "cf2a8d3974ccd2", "6d08f2c539ad01", "Spookify <no-reply@spookify.rdtyads.com>", 2525)

	// PRIMARY
	// Repository
	artistRepository := repository.NewArtistRepositoryImpl(Db)
	tracksRepository := repository.NewTracksRepositoryImpl(Db)
	likesRepository := repository.NewLikeRepositoryImpl(Db)
	usersRepository := repository.NewUserRepositoryImpl(Db)
	tokenRepository := repository.NewTokenRepositoryImpl(Db)
	// Handler
	tracksHandler := handler.NewTracksHandlerImpl(tracksRepository, artistRepository, likesRepository, validators)
	userHandler := handler.NewUserHandlerImpl(validators, usersRepository, tokenRepository, mailer)
	// Router
	router.Init(e, tracksHandler, userHandler)

	// Server (graceful shutdown)
	e.Logger.SetLevel(log.INFO)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// start server
	go func() {
		e.Logger.Printf("Starting server on %s", "8000")
		if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("Server error: %v", err)
		}
	}()

	// wait for interrupt signal to gracefully shutdowns the server with a timeout of 10 seconds.
	<-ctx.Done()
	e.Logger.Print("Shutdown signal received")

	// create a timeout context for shutdown (10 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// attempt gracefully shutdown
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatalf("Server forced to shutdown: %v", err)
	}

	e.Logger.Print("Server gracefully stopped")
}
