package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	mycfg "github.com/sudo-init-do/okies-backend/pkg/config"
	mydb "github.com/sudo-init-do/okies-backend/pkg/db"
	mylog "github.com/sudo-init-do/okies-backend/pkg/logger"
	mw "github.com/sudo-init-do/okies-backend/pkg/middleware"
)

type App struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
	Cfg   *mycfg.Config
}

func main() {
	// global logging setup (zerolog time format etc)
	_ = mylog.New()

	cfg := mycfg.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- Postgres
	pool := mydb.MustOpen(ctx, cfg.DatabaseURL)
	defer pool.Close()

	// --- Redis (optional)
	var rdb *redis.Client
	rc := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	if err := rc.Ping(ctx).Err(); err != nil {
		log.Warn().Err(err).Msg("redis not reachable; proceeding without Redis")
	} else {
		rdb = rc
		defer rdb.Close()
	}

	app := &App{DB: pool, Redis: rdb, Cfg: cfg}

	// --- Router & middleware
	r := chi.NewRouter()
	r.Use(cors.AllowAll().Handler)
	r.Use(mw.Recover)
	r.Use(mw.RequestLog)

	// Health endpoints
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		c, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := app.DB.Ping(c); err != nil {
			http.Error(w, "db_down", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ready"))
	})

	// --- HTTP server with timeouts
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Info().Msgf("API listening on %s (env=%s)", addr, cfg.Env)

	// run server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server crashed")
		}
	}()

	// graceful shutdown
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}
	log.Info().Msg("shutdown complete")
}
