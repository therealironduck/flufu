package orchestrator

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/therealironduck/flufu/internal/agent"
)

func Start(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Go(func() {
		if err := agent.Spawn(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run agent: %v", err)
			os.Exit(1)
		}
	})

	wg.Wait()
}
