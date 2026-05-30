package orchestrator

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/therealironduck/flufu/internal/agent"
	"github.com/therealironduck/flufu/internal/pet"
)

func Start(ctx context.Context) {
	var wg sync.WaitGroup
	cancelCtx, cancel := context.WithCancel(ctx)

	wg.Go(func() {
		if err := agent.Spawn(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run agent: %v", err)
			os.Exit(1)
		}
		cancel()
	})

	wg.Go(func() {
		pet.Render(cancelCtx)
	})

	wg.Wait()
}
