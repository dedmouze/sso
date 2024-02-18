package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sso/internal/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	gw "github.com/dedmouze/protos/gen/go/sso"
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.MustLoad()
	grpcServerEndpoint := flag.String("grpc-server-endpoint", fmt.Sprintf("localhost:%v", cfg.GRPC.Port), "gRPC server endpoint")
	flag.Parse()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gw.RegisterAuthHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	err = gw.RegisterPermissionHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	err = gw.RegisterUserInfoHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.HTTP.Port), mux)
}

func main() {
	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
