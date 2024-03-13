package userInfo

import (
	"context"
	"errors"

	"sso/internal/domain/models"
	"sso/internal/service"

	ssov1 "github.com/dedmouze/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserInfo interface {
	Admin(ctx context.Context, email string) (models.Admin, error)
	User(ctx context.Context, email string) (models.User, error)
}

type serverAPI struct {
	ssov1.UnimplementedUserInfoServer
	userInfo UserInfo
}

func Register(gRPC *grpc.Server, userInfo UserInfo) {
	ssov1.RegisterUserInfoServer(gRPC, &serverAPI{userInfo: userInfo})
}

func (s *serverAPI) User(ctx context.Context, req *ssov1.UserRequest) (*ssov1.UserResponse, error) {
	user, err := s.userInfo.User(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.UserResponse{
		UserID:    user.ID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		VisitedAt: timestamppb.New(user.VisitedAt),
	}, nil
}

func (s *serverAPI) Admin(ctx context.Context, req *ssov1.AdminRequest) (*ssov1.AdminResponse, error) {
	admin, err := s.userInfo.Admin(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, service.ErrAdminNotFound) {
			return nil, status.Error(codes.NotFound, "admin not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.AdminResponse{
		AdminID: admin.ID,
		Email:   admin.Email,
		Level:   int32(admin.Level)}, nil
}
