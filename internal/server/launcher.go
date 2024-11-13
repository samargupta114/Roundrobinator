package server

import (
	"context"
	"sync"

	"github.com/samargupta114/Roundrobinator.git/internal/config"
)

// ServerLauncher defines the interface for launching different server types.
type ServerLauncher interface {
	Launch(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup) error
}
