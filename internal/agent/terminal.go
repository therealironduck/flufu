package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

const (
	bufferSizeInput  int = 256
	bufferSizeOutput int = 4096
)

func Spawn(ctx context.Context) error {
	// Start Claude Code and run it through PTY
	cmd := exec.CommandContext(ctx, "claude")
	cmd.Env = os.Environ()

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("can't start claude code: %w", err)
	}
	defer ptmx.Close()

	// Switch terminal into RAW mode
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("can't switch terminal to raw mode: %w", err)
	}
	defer func() {
		if err := term.Restore(fd, oldState); err != nil {
			fmt.Fprintf(os.Stderr, "error restoring terminal state: %v", err)
		}
	}()

	handleResize(ptmx)
	go handleInput(ptmx)
	go handleOutput(ptmx)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("can't wait for claude execution: %w", err)
	}

	return nil
}

func handleResize(ptmx *os.File) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	go func() {
		for range sigCh {
			setSize(ptmx)
		}
	}()

	setSize(ptmx)
}

func setSize(ptmx *os.File) {
	cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return
	}

	pty.Setsize(ptmx, &pty.Winsize{ //nolint:errcheck
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
}

func handleInput(ptmx *os.File) {
	buf := make([]byte, bufferSizeInput)
	for {
		n, err := os.Stdin.Read(buf)
		if n > 0 {
			_, err = ptmx.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

func handleOutput(ptmx *os.File) {
	buf := make([]byte, bufferSizeOutput)
	for {
		n, err := ptmx.Read(buf)
		if n > 0 {
			_, err = os.Stdout.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}
