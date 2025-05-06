package main

import (
	"context"
	"fmt"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/application/user"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/config"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/db/pgxpool"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/domain"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/infrastructure/db/postgres"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

func main() {
	cfg := config.MustConfig()

	ctx := context.Background()

	pool := pgxpool.MustPGXPool(ctx, cfg, nil, slog.LevelInfo)

	pgRepo := postgres.NewRepository(pool)

	userService := user.NewService(pgRepo)

	if err := userService.CreateUser(ctx, &domain.UserCreate{
		Id:        uuid.New(),
		FirstName: "alex",
		LastName:  nil,
		Email:     "alexxschh1@gmail.com",
		Phone:     nil,
		CreatedAt: time.Now(),
		Password:  "123456",
	}); err != nil {
		panic(err)
	}

	if usr, err := pgRepo.GetUserByEmail(ctx, "alexxschh1@gmail.com"); err != nil {
		panic(err)
	} else {
		fmt.Println("User is", usr)
	}
}
