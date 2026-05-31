package agent

import (
	"fmt"
	"os"
)

func registerHooks() {
	pid := os.Getpid()

	fmt.Printf("%v", pid)
}
