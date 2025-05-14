package app

import (
	"context"
	"fmt"
	userserverpb "github.com/ALexfonSchneider/food-delivery-user-service/gen/grpc/go/user"
	userserver "github.com/ALexfonSchneider/food-delivery-user-service/internal/adapter/grpc/server/user"
	authservice "github.com/ALexfonSchneider/food-delivery-user-service/internal/application/auth"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/application/user"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/config"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/db/pgxpool"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/infrastructure/db/postgres"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/infrastructure/grpc/interceptors"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

func NewApp(ctx context.Context, cfg *config.Config, ready chan<- struct{}) {
	log := slog.Default()

	pool := pgxpool.MustPGXPool(ctx, cfg, nil, slog.LevelInfo)

	pgRepo := postgres.NewRepository(pool)

	if err := pgRepo.Migrate(cfg); err != nil {
		panic(err)
	}

	hasher := authservice.BcryptHasher{}

	jwt := authservice.NewTokenProvider(authservice.Config{
		SecretKey:       cfg.Auth.GetSecretKey(),
		Issuer:          cfg.Auth.GetIssuer(),
		HashAlgorithm:   cfg.Auth.GetHashAlgorithm(),
		AccessTokenTTL:  cfg.Auth.GetAccessTokenTTL(),
		RefreshTokenTTL: cfg.Auth.GetRefreshTokenTTL(),
	})

	userService := user.NewService(pgRepo)
	authService := authservice.NewService(log, pgRepo, hasher, jwt)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.App.GRPCHost, cfg.App.GRPCPort))
	if err != nil {
		panic(err)
	}

	otelOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.NewLoggingInterceptor(log),
			//interceptors.NewAuthInterceptor(log, jwt, userService),
		),
	}

	if tp := otel.GetTracerProvider(); tp != nil {
		otelOpts = append(otelOpts, grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(tp))))
	}

	grpcServer := grpc.NewServer(
		otelOpts...,
	)

	userServer := userserver.NewUserServer(userService, authService)
	userserverpb.RegisterUserServiceServer(grpcServer, userServer)

	go func() {
		log.Info("Starting gRPC server")
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	ready <- struct{}{}

	select {
	case <-ctx.Done():
		log.Info("Shutting down gRPC server")
		grpcServer.GracefulStop()
	}
}
