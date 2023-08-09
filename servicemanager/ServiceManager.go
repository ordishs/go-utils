package servicemanager

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ordishs/go-utils"
	"github.com/ordishs/gocore"
	"golang.org/x/sync/errgroup"
)

type ServiceManager struct {
	services map[string]Service
	logger   utils.Logger
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: make(map[string]Service),
		logger:   gocore.Log("sm"),
	}
}

func (sm *ServiceManager) AddService(name string, service Service) {
	sm.services[name] = service
}

// StartAllAndWait starts all services and waits for them to complete or error.
// If any service errors, all other services are stopped gracefully and the error is returned.
func (sm *ServiceManager) StartAllAndWait(ctx context.Context) error {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Listen for system signals
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		sm.logger.Infof("Received shutdown signal. Stopping services...")
		cancel()
	}()

	g, ctx := errgroup.WithContext(cancelCtx) // Use cancelCtx here

	// Start all services
	for name, service := range sm.services {
		sm.logger.Infof("Starting service %s...", name)

		s := service // capture the loop variable

		g.Go(func() error {
			return s.Start(ctx)
		})
	}

	// Wait for all services to complete or error
	err := g.Wait()
	if err != nil {
		sm.logger.Errorf("Received error: %v", err)
	}

	// Ensure all other services are stopped gracefully with a 10-second timeout
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()

	stopG, _ := errgroup.WithContext(stopCtx) // Use stopCtx here

	for name, service := range sm.services {
		sm.logger.Infof("Stopping service %s...", name)

		s := service // capture the loop variable
		stopG.Go(func() error {
			return s.Stop(stopCtx)
		})
	}

	// Wait for all services to be stopped
	if stopErr := stopG.Wait(); stopErr != nil {
		sm.logger.Warnf("Failed to stop some services: %v", stopErr)
	} else {
		sm.logger.Infof("All services stopped gracefully")
	}

	return err // This is the original error
}
