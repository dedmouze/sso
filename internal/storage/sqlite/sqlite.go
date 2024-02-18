package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new instance of the SQLite storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser saves user to db
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	}

	return id, nil
}

// User returns user model from db by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// IsAdmin returns information whether the user is an admin
func (s *Storage) Admin(ctx context.Context, email string) (models.Admin, error) {
	const op = "storage.sqlite.Admin"

	stmt, err := s.db.Prepare("SELECT id, email, level FROM admins WHERE email = ?")
	if err != nil {
		return models.Admin{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var admin models.Admin
	err = row.Scan(&admin.ID, &admin.Email, &admin.Level)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Admin{}, fmt.Errorf("%s: %w", op, storage.ErrAdminNotFound)
		}
		return models.Admin{}, fmt.Errorf("%s: %w", op, err)
	}

	return admin, nil
}

// App returns app model from db by appID
func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

// AddModerator adds new moderator to admins db
func (s *Storage) AddAdmin(ctx context.Context, email string) error {
	const op = "storage.sqlite.AddAdmin"

	stmt, err := s.db.Prepare("INSERT INTO admins(email, level) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, email, 2)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr); sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrAdminExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// DeleteModerator deletes moderator from admins db by email
func (s *Storage) DeleteAdmin(ctx context.Context, email string) error {
	const op = "storage.sqlite.DeleteAdmin"

	stmt, err := s.db.Prepare("DELETE FROM admins WHERE email = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrAdminNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
