package tty

import (
	"os"
	"sync"
)

var mu sync.Mutex

// Write writes p to stdout atomically with respect to other tty writes.
func Write(p []byte) {
	mu.Lock()
	defer mu.Unlock()

	os.Stdout.Write(p)
}

// WriteString writes s to stdout atomically with respect to other tty writes.
func WriteString(s string) {
	mu.Lock()
	defer mu.Unlock()

	os.Stdout.WriteString(s)
}
