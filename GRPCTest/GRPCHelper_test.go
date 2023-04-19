package greeter

import (
	"context"
	"testing"

	"net"

	"github.com/ordishs/go-utils"
	greeter_api "github.com/ordishs/go-utils/GRPCTest/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GreeterService struct {
	greeter_api.UnimplementedGreeterServiceServer
}

func (s *GreeterService) SayHello(ctx context.Context, req *greeter_api.HelloRequest) (*greeter_api.HelloResponse, error) {
	//return nil, status.Error(codes.Unavailable, "Service is currently unavailable")
	return &greeter_api.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func TestGRPCServerFullCode(t *testing.T) {
	service := &GreeterService{}

	srv := grpc.NewServer()
	greeter_api.RegisterGreeterServiceServer(srv, service)

	go func() {
		// Start the gRPC server
		lis, err := net.Listen("tcp", "localhost:9000")
		assert.NoError(t, err)

		if err := srv.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	// Connect to the server with a gRPC client
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	client := greeter_api.NewGreeterServiceClient(conn)

	// Send a simple request to the server
	req := &greeter_api.HelloRequest{Name: "World"}
	res, err := client.SayHello(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "Hello, World", res.Message)
}

func TestGRPCServerUsingGRPCHelper(t *testing.T) {
	srv, err := utils.GetGRPCServer(&utils.ConnectionOptions{
		Tracer: true,
	})
	require.NoError(t, err)

	service := &GreeterService{}
	greeter_api.RegisterGreeterServiceServer(srv, service)

	go func() {
		// Start the gRPC server
		lis, err := net.Listen("tcp", "localhost:9000")
		assert.NoError(t, err)

		if err := srv.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	// Connect to the server with a gRPC client
	conn, err := utils.GetGRPCClient(context.Background(), "localhost:9000", &utils.ConnectionOptions{
		Tracer:     true,
		MaxRetries: 10,
	})
	require.NoError(t, err)

	defer conn.Close()

	client := greeter_api.NewGreeterServiceClient(conn)

	// Send a simple request to the server
	req := &greeter_api.HelloRequest{Name: "World"}
	res, err := client.SayHello(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "Hello, World", res.Message)
}
