package servicemanager

import "context"

// Service is an interface for services that can be started and stopped.
type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
