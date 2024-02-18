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

type Key string

var validationErrorKey = Key("validation error")

func UnaryValidationInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		const op = "grpc.interceptor.UnaryValidationInterceptor"

		method := info.FullMethod
		log := log.With(slog.String("op", op), slog.String("method", method))

		log.Info("validation interceptor enabled")

		var err error

		switch method {
		case "/auth.Auth/Login":
			err = validateLogin(req.(*ssov1.LoginRequest))
		case "/auth.Auth/Register":
			err = validateRegister(req.(*ssov1.RegisterRequest))
		case "/userInfo.UserInfo/User":
			err = validateEmailToken(req.(*ssov1.UserRequest))
		case "/userInfo.UserInfo/Admin":
			err = validateEmailToken(req.(*ssov1.AdminRequest))
		case "/permission.Permission/AddAdmin":
			err = validateEmailToken(req.(*ssov1.AddAdminRequest))
		case "/permission.Permission/DeleteAdmin":
			err = validateEmailToken(req.(*ssov1.DeleteAdminRequest))
		default:
			err = status.Error(codes.Unimplemented, "method not found")
		}

		if err != nil {
			log.Error("validation error", sl.Err(err))
			ctx = context.WithValue(ctx, validationErrorKey, err)
		} else {
			log.Info("request validated")
		}

		return handler(ctx, req)
	}
}

func ValidateRequest(ctx context.Context) error {
	validationError, ok := validationErrorFromContext(ctx)
	if !ok && validationError != nil {
		return status.Error(codes.Internal, "internal error")
	}
	return validationError
}

func validationErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(validationErrorKey).(error)
	return err, ok
}

const emptyValue = 0

var (
	emailRequired    = "email is required"
	passwordRequired = "password is required"
	tokenRequired    = "token is required"
	appIDRequired    = "app_id is required"
)

type requestTokenEmail interface {
	GetEmail() string
	GetToken() string
}

func validateEmailToken(req requestTokenEmail) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, emailRequired)
	}
	if req.GetToken() == "" {
		return status.Error(codes.InvalidArgument, tokenRequired)
	}
	return nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, appIDRequired)
	}
	return validateEmailPassword(req.GetEmail(), req.GetPassword())
}

func validateRegister(req *ssov1.RegisterRequest) error {
	return validateEmailPassword(req.GetEmail(), req.GetPassword())
}

func validateEmailPassword(email, password string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, emailRequired)
	}
	if password == "" {
		return status.Error(codes.InvalidArgument, passwordRequired)
	}
	return nil
}
