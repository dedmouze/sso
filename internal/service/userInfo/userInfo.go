package userInfo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"sso/internal/domain/models"
	"sso/internal/lib/logger/sl"
	"sso/internal/service"
	"sso/internal/storage"
)

type UserInfo struct {
	log          *slog.Logger
	userProvider UserProvider
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	Admin(ctx context.Context, email string) (models.Admin, error)
}

// New returns a new instance of the UserInfo service
func New(
	log *slog.Logger,
	userProvider UserProvider,
) *UserInfo {
	return &UserInfo{
		log:          log,
		userProvider: userProvider,
	}
}

// Admin checks if user is admin
//
// If user with given userID doesn't exist, returns error
func (u *UserInfo) Admin(
	ctx context.Context,
	email string,
) (models.Admin, error) {
	const op = "services.userInfo.Admin"

	log := u.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("checking if user is admin")

	admin, err := u.userProvider.Admin(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrAdminNotFound) {
			log.Warn("admin not found", sl.Err(err))
			return models.Admin{}, fmt.Errorf("%s: %w", op, service.ErrAdminNotFound)
		}

		log.Error("failed to check if user is admin")
		return models.Admin{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", admin.Level > 1))

	return admin, nil
}

// User returns user info
//
// If user with given email doesn't exist, returns error
func (u *UserInfo) User(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "services.userInfo.User"

	log := u.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("getting user info")

	user, err := u.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, service.ErrUserNotFound)
		}

		log.Error("failed to get user info")
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user found")

	return user, nil
}
