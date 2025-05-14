package main

import (
	"context"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/app"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/config"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/infrastructure/otel"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustConfig()

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	_, shutdownOtel := otel.InitTracerProvider(cfg)

	ready := make(chan<- struct{}, 1)
	app.NewApp(ctx, cfg, ready)

	slog.Info("shutdown otel")
	shutdownOtel()
}
