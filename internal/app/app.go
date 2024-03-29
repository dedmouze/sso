package app

import (
	"log/slog"
	"time"

	"sso/internal/app/grpcapp"
	"sso/internal/service/auth"
	"sso/internal/service/permission"
	"sso/internal/service/userInfo"
	"sso/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
	userKey string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, storage, storage, tokenTTL, userKey)
	userInfoService := userInfo.New(log, storage)
	permissionService := permission.New(log, storage, storage)

	grpcApp := grpcapp.New(log, authService, userInfoService, permissionService, storage, grpcPort, userKey)
	return &App{
		GRPCServer: grpcApp,
	}
}
