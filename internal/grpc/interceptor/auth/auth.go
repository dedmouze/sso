package auth

import (
	"context"
	"errors"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	"sso/internal/storage"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	permissionRequired = []string{
		"/userInfo.UserInfo/User",
		"/userInfo.UserInfo/Admin",
		"/permission.Permission/AddAdmin",
		"/permission.Permission/DeleteAdmin",
	}
)

type AppProvider interface {
	AppByKey(ctx context.Context, apiKey string) (models.App, error)
}

func UnaryAuthenticationInterceptor(log *slog.Logger, appProvider AppProvider, userKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		const op = "grpc.interceptor.UnaryAuthenticationInterceptor"

		method := info.FullMethod
		log := log.With(slog.String("op", op), slog.String("method", method))

		log.Info("auth interceptor enabled")

		var err error

		for _, s := range permissionRequired {
			if s == method {
				md, ok := metadata.FromIncomingContext(ctx)
				if !ok {
					err = status.Error(codes.InvalidArgument, "missing metadata")
					break
				}
				err = valid(md["authorization"], log, appProvider, userKey)
				break
			}
		}

		if err != nil {
			log.Warn("auth error", sl.Err(err))
			return nil, err
		}

		log.Info("request authenticated")

		return handler(ctx, req)
	}
}

var (
	internalErr = status.Error(codes.Internal, "internal error")
	authErr     = status.Error(codes.Unauthenticated, "invalid token or key")
)

func valid(authorization []string, log *slog.Logger, appProvider AppProvider, userKey string) error {
	const op = "interceptor.auth.valid"

	log = log.With("op", op)

	if len(authorization) < 1 {
		return authErr
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")

	isUser := false
	for _, c := range token {
		if c == rune('.') {
			isUser = true
			break
		}
	}

	if isUser {
		token, err := jwt.Parse(token, userKey)
		if err != nil {
			return internalErr
		}

		if token.Level < 2 {
			return authErr
		}

		log.Info("user is authenticated")
	} else {
		_, err := appProvider.AppByKey(context.Background(), token)
		if err != nil {
			if errors.Is(err, storage.ErrAppNotFound) {
				return authErr
			}
			return internalErr
		}
		log.Info("app is authenticated")
	}

	return nil
}
