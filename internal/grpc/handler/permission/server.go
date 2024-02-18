package permission

import (
	"context"
	"errors"
	"sso/internal/grpc/interceptor/auth"
	"sso/internal/grpc/interceptor/validation"
	"sso/internal/service"

	ssov1 "github.com/dedmouze/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Permission interface {
	AddAdmin(ctx context.Context, email string) error
	DeleteAdmin(ctx context.Context, email string) error
}

type serverAPI struct {
	ssov1.UnimplementedPermissionServer
	permission Permission
}

func Register(gRPC *grpc.Server, permission Permission) {
	ssov1.RegisterPermissionServer(gRPC, &serverAPI{permission: permission})
}

func (s *serverAPI) AddAdmin(ctx context.Context, req *ssov1.AddAdminRequest) (*ssov1.AddAdminResponse, error) {
	if err := validation.ValidateRequest(ctx); err != nil {
		return nil, err
	}

	if err := auth.CheckRequest(ctx); err != nil {
		return nil, err
	}

	err := s.permission.AddAdmin(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, service.ErrAdminExists) {
			return nil, status.Error(codes.AlreadyExists, "admin already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.AddAdminResponse{}, nil
}

func (s *serverAPI) DeleteAdmin(ctx context.Context, req *ssov1.DeleteAdminRequest) (*ssov1.DeleteAdminResponse, error) {
	if err := validation.ValidateRequest(ctx); err != nil {
		return nil, err
	}

	if err := auth.CheckRequest(ctx); err != nil {
		return nil, err
	}

	err := s.permission.DeleteAdmin(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, service.ErrAdminNotFound) {
			return nil, status.Error(codes.NotFound, "admin not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.DeleteAdminResponse{}, nil
}
