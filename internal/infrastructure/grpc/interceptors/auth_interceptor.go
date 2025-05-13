package interceptors

import (
	"context"
	"fmt"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/application/auth"
	"github.com/ALexfonSchneider/food-delivery-user-service/internal/application/user"
	"github.com/ALexfonSchneider/food-delivery-user-service/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var publicMethods = map[string]bool{
	"/user.UserService/LoginUser":    true,
	"/user.UserService/RegisterUser": true,
}

func NewAuthInterceptor(log *slog.Logger, jwt *auth.TokenProvider, userService *user.Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		log = log.With("method", info.FullMethod)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Error("missing metadata")
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}

		authHeaders := md["authorization"]
		if len(authHeaders) == 0 {
			log.Error("missing authorization header")
			return nil, status.Error(codes.InvalidArgument, "missing authorization header")
		}

		authHeader := authHeaders[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, fmt.Errorf("invalid authorization format")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateToken(token)
		if err != nil {
			log.Error("invalid token", logger.Err(err))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		usr, err := userService.GetUserById(ctx, claims.UserCredentials.UserID)
		if err != nil {
			log.Error("failed to get user", logger.Err(err))
			return nil, status.Error(codes.Unauthenticated, "invalid user")
		}

		ctx = context.WithValue(ctx, "user", usr)

		return handler(ctx, req)
	}
}
