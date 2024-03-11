package permission

import (
	"context"
	"errors"
	"sso/internal/service"

	ssov1 "github.com/dedmouze/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *serverAPI) AddAdmin(ctx context.Context, req *ssov1.AddAdminRequest) (*emptypb.Empty, error) {
	err := s.permission.AddAdmin(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, service.ErrAdminExists) {
			return nil, status.Error(codes.AlreadyExists, "admin already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}

func (s *serverAPI) DeleteAdmin(ctx context.Context, req *ssov1.DeleteAdminRequest) (*emptypb.Empty, error) {
	err := s.permission.DeleteAdmin(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, service.ErrAdminNotFound) {
			return nil, status.Error(codes.NotFound, "admin not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &emptypb.Empty{}, nil
}
