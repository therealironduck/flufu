package pet

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	hookLogPath = "/tmp/flufu-hook.log"
	hookLogPerm = 0o644
)

func getSocketForPid(pid int) string {
	return fmt.Sprintf("/tmp/flufu-%d.sock", pid)
}

func getSocket() string {
	return getSocketForPid(os.Getpid())
}

func Listen(ctx context.Context) error {
	lc := net.ListenConfig{}

	ln, err := lc.Listen(ctx, "unix", getSocket())
	if err != nil {
		return fmt.Errorf("cant listen for unix socket: %w", err)
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	con, err := ln.Accept()
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}

		return fmt.Errorf("cant accept unix socket connection: %w", err)
	}
	defer con.Close()

	msg, err := io.ReadAll(con)
	if err != nil {
		return fmt.Errorf("cant read from unix socket: %w", err)
	}

	_ = os.WriteFile(hookLogPath, append(msg, '\n'), hookLogPerm)

	return nil
}

func Send(ctx context.Context, pid int, message string) error {
	d := net.Dialer{}
	con, err := d.DialContext(ctx, "unix", getSocketForPid(pid))
	if err != nil {
		return fmt.Errorf("cant connect to flufu socket: %w", err)
	}
	defer con.Close()

	if _, err = fmt.Fprint(con, message); err != nil {
		return fmt.Errorf("cant write to flufu socket: %w", err)
	}

	return nil
}
