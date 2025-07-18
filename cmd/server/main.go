package main

import (
	"context"
	"fmt"
	"gorestsubs/internal/handler"
	"gorestsubs/internal/repository"
	"gorestsubs/internal/service"
	"gorestsubs/pkg/config"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gorestsubs/docs"
)

func main() {
	cfg := config.MustLoad()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		logger.Error("failed to create migrate instance", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("failed to apply migrations", "error", err)
		os.Exit(1)
	}

	repo := repository.NewSubscriptionRepository(pool)
	svc := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(svc)

	router := gin.Default()
	h.Register(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("starting server", "port", cfg.HTTP.Port)

	if err := http.ListenAndServe(":"+cfg.HTTP.Port, router); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
