package auth

import (
	"context"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"

	ssov1 "github.com/dedmouze/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppGetter interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type Key string

// internal services
var (
	permission = 2
	userInfo   = 3
)

var authErrorKey = Key("auth error")

func UnaryAuthenticationInterceptor(log *slog.Logger, appGetter AppGetter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		const op = "grpc.interceptor.UnaryAuthenticationInterceptor"

		log := log.With(slog.String("op", op))

		log.Info("auth interceptor enabled")

		var err error
		method := info.FullMethod

		switch method {
		case "/userInfo.UserInfo/User":
			err = authenticate(ctx, req.(*ssov1.UserRequest), appGetter, userInfo)
		case "/userInfo.UserInfo/Admin":
			err = authenticate(ctx, req.(*ssov1.AdminRequest), appGetter, userInfo)
		case "/permission.Permission.AddAdmin":
			err = authenticate(ctx, req.(*ssov1.AddAdminRequest), appGetter, permission)
		case "/permission.Permission.DeleteAdmin":
			err = authenticate(ctx, req.(*ssov1.DeleteAdminRequest), appGetter, permission)
		}

		if err != nil {
			log.Error("auth error", sl.Err(err))
			ctx = context.WithValue(ctx, authErrorKey, err)
		} else {
			log.Info("request authenticated")
		}

		return handler(ctx, req)
	}
}

func CheckRequest(ctx context.Context) error {
	authError, ok := authErrorFromContext(ctx)
	if !ok && authError != nil {
		return status.Error(codes.Internal, "internal error")
	}
	return authError
}

func authErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(authErrorKey).(error)
	return err, ok
}

type requestToken interface {
	GetToken() string
}

func authenticate(
	ctx context.Context,
	req requestToken,
	appGetter AppGetter,
	appID int,
) error {
	app, err := appGetter.App(ctx, appID)
	if err != nil {
		return status.Error(codes.Internal, "internal error")
	}

	raw := req.GetToken()
	token, err := jwt.Parse(raw, app.Secret)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if token.Level < 2 {
		return status.Error(codes.PermissionDenied, "dont have permission")
	}

	return nil
}
