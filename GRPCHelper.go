package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------
// The PasswordCredentials type and the receivers it implements, allow
// us to use the grpc.WithPerRPCCredentials() dial option to pass
// credentials to downstream middleware
type PasswordCredentials map[string]string

func NewPassCredentials(m map[string]string) PasswordCredentials {
	return PasswordCredentials(m)
}

func (pc PasswordCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return pc, nil
}

func (PasswordCredentials) RequireTransportSecurity() bool {
	return false
}

// ---------------------------------------------------------------------

type ConnectionOptions struct {
	MaxMessageSize int                 // Max message size in bytes
	SecurityLevel  int                 // 0 = insecure, 1 = secure, 2 = secure with client cert
	Tracer         bool                // Enable OpenTelemetry tracing
	CertFile       string              // CA cert file if SecurityLevel > 0
	CaCertFile     string              // CA cert file if SecurityLevel > 0
	KeyFile        string              // Client key file if SecurityLevel > 1
	MaxRetries     int                 // Max number of retries for transient errors
	RetryBackoff   time.Duration       // Backoff between retries
	Credentials    PasswordCredentials // Credentials to pass to downstream middleware (optional)
}

// ---------------------------------------------------------------------

func init() {
	// The secret sauce
	resolver.SetDefaultScheme("dns")
}

func GetGRPCClient(ctx context.Context, address string, connectionOptions *ConnectionOptions) (*grpc.ClientConn, error) {
	if address == "" {
		return nil, errors.New("address is required")
	}

	if err := validateConnectionOptions(connectionOptions); err != nil {
		return nil, err
	}

	if connectionOptions.MaxMessageSize == 0 {
		connectionOptions.MaxMessageSize = 1024 * 1024 * 1000 // 1GB: transactions are getting bigger, current limit is 9MB
	}

	opts := []grpc.DialOption{
		// grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(connectionOptions.MaxMessageSize),
		),
	}

	if connectionOptions.SecurityLevel == 0 {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		caCert, err := os.ReadFile(connectionOptions.CaCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read ca cert file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		cert, err := tls.LoadX509KeyPair(connectionOptions.CertFile, connectionOptions.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load cert key pair: %w", err)
		}

		config := &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: false,
			Certificates:       []tls.Certificate{cert},
		}

		opts = append(
			opts,
			grpc.WithTransportCredentials(credentials.NewTLS(config)),
		)
	}

	if connectionOptions.Tracer {
		opts = append(
			opts,
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		)
	}

	if connectionOptions.Credentials != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(connectionOptions.Credentials))
	}

	// Retry interceptor...
	if connectionOptions.MaxRetries > 0 {
		if connectionOptions.RetryBackoff == 0 {
			connectionOptions.RetryBackoff = 100 * time.Millisecond
		}

		opts = append(opts, grpc.WithUnaryInterceptor(retryInterceptor(connectionOptions.MaxRetries, connectionOptions.RetryBackoff)))
	}

	conn, err := grpc.DialContext(
		ctx,
		address,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("Error dialling grpc service at %s: %v", address, err)
	}

	return conn, nil
}

func GetGRPCServer(connectionOptions *ConnectionOptions) (*grpc.Server, error) {
	if err := validateConnectionOptions(connectionOptions); err != nil {
		return nil, err
	}

	var opts []grpc.ServerOption

	if connectionOptions.MaxMessageSize == 0 {
		connectionOptions.MaxMessageSize = 1024 * 1024 * 1000 // 1GB: transactions are getting bigger, current limit is 9MB
	}

	opts = append(opts, grpc.MaxRecvMsgSize(connectionOptions.MaxMessageSize))

	if connectionOptions.Tracer {
		opts = append(opts, grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
		opts = append(opts, grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
	}

	if connectionOptions.SecurityLevel > 0 {
		tlsConfig := &tls.Config{}
		var err error

		// For LEVEL1 and LEVEL2 we need to encrypt all traffic.  Client will only verify our certificate in LEVEL2

		// Create a CA certificate pool and add ca.crt to it
		caCert, err := os.ReadFile(connectionOptions.CertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read ca cert file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.RootCAs = caCertPool
		tlsConfig.ClientCAs = caCertPool

		cert, err := tls.LoadX509KeyPair(connectionOptions.CertFile, connectionOptions.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read key pair: %w", err)
		}

		tlsConfig.Certificates = []tls.Certificate{cert}

		if connectionOptions.SecurityLevel == 2 {
			tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
		} else if connectionOptions.SecurityLevel == 3 {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}

		creds := credentials.NewTLS(tlsConfig)

		opts = append(opts, grpc.Creds(creds))
	}

	return grpc.NewServer(opts...), nil
}

func validateConnectionOptions(connectionData *ConnectionOptions) error {
	if connectionData.SecurityLevel > 0 {
		if connectionData.CertFile == "" {
			return errors.New("certFile is required for a secure connection")
		}

		if connectionData.CaCertFile == "" {
			return errors.New("caCertFile is required for a secure connection")
		}

		if connectionData.SecurityLevel > 1 {
			if connectionData.KeyFile == "" {
				return errors.New("keyFile is required for a secure connection")
			}
		}
	}

	return nil
}

func retryInterceptor(maxRetries int, retryBackoff time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var err error

		for i := 0; i < maxRetries; i++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}

			// Check if we can retry (e.g., codes.Unavailable, codes.DeadlineExceeded)
			if status.Code(err) != codes.Unavailable && status.Code(err) != codes.DeadlineExceeded {
				break
			}

			log.Printf("Retry attempt %d for request: %s\n", i+1, method)
			time.Sleep(retryBackoff)
		}

		return err
	}
}
