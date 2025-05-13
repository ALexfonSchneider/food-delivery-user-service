package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

func NewLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		log.Info("Incoming request",
			"method", info.FullMethod,
			"type", "unary",
		)

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		level := slog.LevelInfo
		if err != nil {
			level = slog.LevelError
		}

		log.LogAttrs(ctx, level,
			"Finished call",
			slog.String("method", info.FullMethod),
			slog.Duration("duration", duration),
			slog.Any("error", err),
		)

		return resp, err
	}
}
