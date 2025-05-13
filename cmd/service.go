package main

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
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/infrastructure/otel"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"sync"
)

func main() {
	cfg := config.MustConfig()

	ctx := context.Background()

	tp, _ := otel.InitTracer(cfg)

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

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", cfg.App.GRPCPort))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.NewLoggingInterceptor(log),
			//interceptors.NewAuthInterceptor(log, jwt, userService),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(tp))),
	)

	userServer := userserver.NewUserServer(userService, authService)
	userserverpb.RegisterUserServiceServer(grpcServer, userServer)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Info("Starting gRPC server")
		if err = grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
