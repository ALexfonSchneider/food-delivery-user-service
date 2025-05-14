package integrational

import (
	"context"
	"fmt"
	"github.com/ALexfonSchneider/food-delivery-user-service/gen/grpc/go/user"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/app"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/config"
	"github.com/ALexfonSchneider/food-delivery-user-service/pkg/random"
	"google.golang.org/grpc"
	"log"
	"sync"
	"testing"
)

func runClient(ctx context.Context, wg *sync.WaitGroup, client user.UserServiceClient) {
	defer wg.Done()

	request := &user.RegisterUserRequest{
		FirstName: random.RandStringRunes(10),
		LastName:  nil,
		Email:     fmt.Sprintf("%v@example.com", random.RandStringRunes(10)),
		Phone:     random.RandStringRunes(10),
		Password:  random.RandStringRunes(10),
	}

	_, err := client.RegisterUser(ctx, request)
	if err != nil {
		log.Fatalf("%v.RegisterUser(_) = _, %v", client, err)
	}
}

func TestLogin_Massive(t *testing.T) {
	cfg := config.MustConfig()

	ctx, cancel := context.WithCancel(context.Background())

	ready := make(chan struct{}, 1)
	go app.NewApp(ctx, cfg, ready)
	<-ready

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.App.GRPCHost, cfg.App.GRPCPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)

	wg := &sync.WaitGroup{}
	numClients := 1000 // число одновременных клиентов
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go runClient(ctx, wg, client)
	}
	wg.Wait()

	cancel()
}
