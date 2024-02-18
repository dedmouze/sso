package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"sso/internal/grpc/handler/auth"
	"sso/internal/grpc/handler/permission"
	"sso/internal/grpc/handler/userInfo"
	authInterceptor "sso/internal/grpc/interceptor/auth"
	"sso/internal/grpc/interceptor/validation"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService auth.Auth,
	userInfoService userInfo.UserInfo,
	permissionService permission.Permission,
	appGetter authInterceptor.AppGetter,
	port int,
) *App {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			validation.UnaryValidationInterceptor(log),
			authInterceptor.UnaryAuthenticationInterceptor(log, appGetter),
		),
	)

	// TODO: remove
	// gRPCServer := grpc.NewServer(
	// 	grpc.UnaryInterceptor(validation.UnaryValidationInterceptor(log)),
	// )

	auth.Register(gRPCServer, authService)
	userInfo.Register(gRPCServer, userInfoService)
	permission.Register(gRPCServer, permissionService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting gRPC server", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
