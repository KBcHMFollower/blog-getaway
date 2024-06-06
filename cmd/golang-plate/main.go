package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-plate/internal/config"
	"test-plate/internal/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	cfg := config.MustLoad()

	log:= logger.CongigurateLogger(cfg.Env)

	log.Debug("debug messages are enabled")

	router := chi.NewRouter();

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/auth", func (r chi.Router){
		r.Get("/", func (w http.ResponseWriter, r *http.Request){
			render.JSON(w, r, "you tru get auth")
		})
	})

	log.Info(
		"Get Away Api try starting with", 
		slog.String("env", cfg.Env),
		slog.String("url", cfg.Address)	,
	)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout: cfg.IddleTimeout,
	}

	go func(){
		if err := srv.ListenAndServe(); err != nil{
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", err)

		return
	}

	log.Info("server stopped")
}