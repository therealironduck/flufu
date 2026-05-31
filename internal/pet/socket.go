package pet

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
)


func getSocketForPid(pid int) string {
	return fmt.Sprintf("/tmp/flufu-%d.sock", pid)
}

func getSocket() string {
	return getSocketForPid(os.Getpid())
}

func Listen(ctx context.Context, msgCh chan<- string) error {
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

	for {
		con, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf("cant accept unix socket connection: %w", err)
		}

		go func() {
			defer con.Close()

			msg, err := io.ReadAll(con)
			if err != nil {
				return
			}

			select {
			case msgCh <- string(msg):
			default:
			}
		}()
	}
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
