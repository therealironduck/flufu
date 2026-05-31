package orchestrator

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/therealironduck/flufu/internal/agent"
	"github.com/therealironduck/flufu/internal/ai"
	"github.com/therealironduck/flufu/internal/pet"
)

func Start(ctx context.Context) {
	var wg sync.WaitGroup
	cancelCtx, cancel := context.WithCancel(ctx)

	aiInstance := ai.New()

	wg.Go(func() {
		if err := aiInstance.Init(cancelCtx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to init AI: %v", err)
		}
	})

	msgCh := make(chan string, 1)

	wg.Go(func() {
		if err := agent.Spawn(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run agent: %v", err)
			os.Exit(1)
		}
		cancel()
	})

	wg.Go(func() {
		pet.Render(cancelCtx, msgCh)
	})

	wg.Go(func() {
		if err := pet.Listen(cancelCtx, msgCh); err != nil {
			fmt.Fprintf(os.Stderr, "failed to listen for socket: %v", err)
			os.Exit(1)
		}
	})

	wg.Wait()
}
