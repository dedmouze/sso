package validation

import (
	"context"
	"log/slog"
	"sso/internal/lib/logger/sl"

	ssov1 "github.com/dedmouze/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryValidationInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		const op = "grpc.interceptor.UnaryValidationInterceptor"

		method := info.FullMethod
		log := log.With(slog.String("op", op), slog.String("method", method))

		log.Info("validation interceptor enabled")

		var err error

		switch method {
		case "/auth.Auth/Login":
			err = validateEmailPassword(req.(*ssov1.LoginRequest))
		case "/auth.Auth/Register":
			err = validateEmailPassword(req.(*ssov1.RegisterRequest))
		case "/auth.Auth/RegisterApp":
			err = validateRegisterApp(req.(*ssov1.RegisterAppRequest))
		case "/userInfo.UserInfo/User":
			err = validateEmail(req.(*ssov1.UserRequest))
		case "/userInfo.UserInfo/Admin":
			err = validateEmail(req.(*ssov1.AdminRequest))
		case "/permission.Permission/AddAdmin":
			err = validateEmail(req.(*ssov1.AddAdminRequest))
		case "/permission.Permission/DeleteAdmin":
			err = validateEmail(req.(*ssov1.DeleteAdminRequest))
		default:
			err = status.Error(codes.Unimplemented, "method not found")
		}

		if err != nil {
			log.Warn("validation error", sl.Err(err))
			return nil, err
		}
		log.Info("request validated")

		return handler(ctx, req)
	}
}

var (
	emailRequired    = "email is required"
	nameRequired     = "name is required"
	passwordRequired = "password is required"
)

type requestEmail interface {
	GetEmail() string
}

type requestEmailPassword interface {
	requestEmail
	GetPassword() string
}

func validateEmail(req requestEmail) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, emailRequired)
	}
	return nil
}

func validateEmailPassword(req requestEmailPassword) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, emailRequired)
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, passwordRequired)
	}
	return nil
}

func validateRegisterApp(req *ssov1.RegisterAppRequest) error {
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, nameRequired)
	}
	return nil
}
