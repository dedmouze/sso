package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/sl"
	"sso/internal/lib/secret"
	"sso/internal/service"
	"sso/internal/service/userInfo"
	"sso/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider userInfo.UserProvider
	userChanger  UserChanger
	appProvider  AppProvider
	appSaver     AppSaver
	tokenTTL     time.Duration
	userKey      string
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (userID int64, err error)
}

type UserChanger interface {
	UpdateUserVisitTime(
		ctx context.Context,
		email string,
		visitTime time.Time,
	) error
}

type AppSaver interface {
	SaveApp(
		ctx context.Context,
		name string,
		apiKey string,
	) error
}

type AppProvider interface {
	AppByID(ctx context.Context, appID int) (models.App, error)
	AppByKey(ctx context.Context, apiKey string) (models.App, error)
}

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider userInfo.UserProvider,
	userChanger UserChanger,
	appSaver AppSaver,
	appProvider AppProvider,
	tokenTTL time.Duration,
	userKey string,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		userChanger:  userChanger,
		appSaver:     appSaver,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
		userKey:      userKey,
	}
}

// Login checks if user with given credentials exists in system
//
// If user exists, but password is incorrect, returns error
// If user doesn't exist, returns error
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	const op = "services.auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email), //optional
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, service.ErrInvalidCredentials)
	}

	admin, err := a.userProvider.Admin(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrAdminNotFound) {
			log.Info("user not admin")
		} else {
			log.Error("failed to get admin", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("user logged in successfully")

	err = a.userChanger.UpdateUserVisitTime(ctx, email, time.Now())
	if err != nil {
		log.Warn("failed to update visit time")
	}

	if err != nil {
		admin.Level = 1 // TODO: remove this brute force approach
	}
	token, err := jwt.NewToken(user, admin, a.tokenTTL, a.userKey)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("token genereted")

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID
//
// If user with given username already exists, returns error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "services.auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email), //optional
	)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, service.ErrUserExists)
		}

		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

// RegisterNewApp registers new app in the system and returns apiKey
//
// If app with given name already exists, returns error
func (a *Auth) RegisterNewApp(ctx context.Context, name string) (string, string, error) {
	const op = "service.auth.RegisterNewApp"

	log := a.log.With(
		slog.String("op", op),
		slog.String("app name", name),
	)

	log.Info("registering new app")

	var apiKey string
	var err error
	for {
		apiKey, err = secret.GenerateSecret()
		if err != nil {
			log.Error("failed to generate secret")
			return "", "", fmt.Errorf("%s: %w", op, err)
		}

		_, err := a.appProvider.AppByKey(ctx, apiKey)
		if err != nil {
			if errors.Is(err, storage.ErrAppNotFound) {
				break
			}
			log.Error("failed to get app by key", sl.Err(err))
			return "", "", fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := a.appSaver.SaveApp(ctx, name, apiKey); err != nil {
		if errors.Is(err, storage.ErrAppExists) {
			log.Warn("app already exists")
			return "", "", fmt.Errorf("%s: %w", op, service.ErrAppAlreadyExists)
		}
		log.Error("failed to save app", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("app registered")

	return apiKey, a.userKey, nil
}
