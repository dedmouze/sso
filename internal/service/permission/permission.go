package permission

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"sso/internal/lib/logger/sl"
	"sso/internal/service"
	"sso/internal/storage"
)

type Permission struct {
	log          *slog.Logger
	adminAdder   AdminAdder
	adminDeleter AdminDeleter
}

type AdminAdder interface {
	AddAdmin(ctx context.Context, email string) error
}

type AdminDeleter interface {
	DeleteAdmin(ctx context.Context, email string) error
}

func New(
	log *slog.Logger,
	adminAdder AdminAdder,
	adminDeleter AdminDeleter,
) *Permission {
	return &Permission{
		log:          log,
		adminAdder:   adminAdder,
		adminDeleter: adminDeleter,
	}
}

func (p *Permission) AddAdmin(ctx context.Context, email string) error {
	const op = "services.permission.AddAdmin"

	log := p.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("adding admin")

	if err := p.adminAdder.AddAdmin(ctx, email); err != nil {
		if errors.Is(err, storage.ErrAdminExists) {
			log.Warn("admin already exists", sl.Err(err))
			return fmt.Errorf("%s: %w", op, service.ErrAdminExists)
		}

		log.Error("failed to add admin")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("admin added successfully")

	return nil
}

func (p *Permission) DeleteAdmin(ctx context.Context, email string) error {
	const op = "services.permission.DeleteAdmin"

	log := p.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("deleting admin")

	if err := p.adminDeleter.DeleteAdmin(ctx, email); err != nil {
		if errors.Is(err, storage.ErrAdminNotFound) {
			log.Warn("admin not found", sl.Err(err))
			return fmt.Errorf("%s: %w", op, service.ErrAdminNotFound)
		}

		log.Error("failed to delete admin")
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("admin deleted")

	return nil
}
