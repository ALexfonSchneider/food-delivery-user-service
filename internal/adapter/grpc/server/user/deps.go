package user

import (
	"github.com/ALexfonSchneider/food-delivery-user-service/gen/grpc/go/user"
	authservice "github.com/ALexfonSchneider/food-delivery-user-service/internal/application/auth"
	userservice "github.com/ALexfonSchneider/food-delivery-user-service/internal/application/user"
)

type Server struct {
	user.UnimplementedUserServiceServer
	user *userservice.Service
	auth *authservice.Service
}

func NewUserServer(userService *userservice.Service, authService *authservice.Service) *Server {
	return &Server{user: userService, auth: authService}
}
