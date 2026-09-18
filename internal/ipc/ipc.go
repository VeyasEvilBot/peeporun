// Package ipc lets a short-lived `peeporun <command>` invocation control an
// already-running peepoRun TUI instance, over a local Unix domain socket.
//
// This is what makes external hotkey binding (e.g. KDE Custom Shortcuts
// running `peeporun split`) possible without any OS-level global key
// capture - the desktop environment's own trusted shortcut system runs a
// normal command, which just talks to the real running instance over this
// socket. Works identically under X11 and Wayland, since it never touches
// input devices or window-system APIs at all.
package ipc

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Command is one request from a client (`peeporun hit`, etc.) to the
// running TUI instance. Action is one of "hit", "undo", "split", "reset",
// "preset" (Arg is the preset ID for that last one). Result must be sent
// on exactly once by whoever handles the command.
type Command struct {
	Action string
	Arg    string
	Result chan Result
}

// Result is the outcome of handling a Command, sent back to the waiting
// client connection.
type Result struct {
	OK  bool
	Msg string
}

// Serve starts listening on socketPath and, for every line received on a
// connection, builds a Command (with a fresh Result channel) and passes it
// to handle. handle must eventually send exactly one Result on cmd.Result -
// typically by routing the Command into a tea.Program via p.Send(cmd) and
// having the Bubble Tea Update() loop reply on the channel once it's
// processed the action, since that's the only goroutine allowed to touch
// the app's state.
//
// Returns a stop func to shut the listener down and remove the socket
// file - call it on program exit.
func Serve(socketPath string, handle func(Command)) (stop func(), err error) {
	// Best-effort: clean up a stale socket file left behind by a previous
	// run that didn't exit cleanly (e.g. a crash). If another instance is
	// genuinely still running, its Listen below will simply fail instead,
	// which is the correct outcome.
	os.Remove(socketPath)

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", socketPath, err)
	}

	done := make(chan struct{})
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-done:
					return // listener closed intentionally by stop()
				default:
					continue
				}
			}
			go serveConn(conn, handle)
		}
	}()

	stop = func() {
		close(done)
		ln.Close()
		os.Remove(socketPath)
	}
	return stop, nil
}

func serveConn(conn net.Conn, handle func(Command)) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	action, arg, _ := strings.Cut(strings.TrimSpace(scanner.Text()), " ")
	if action == "" {
		fmt.Fprintln(conn, "ERR empty command")
		return
	}

	resultCh := make(chan Result, 1)
	handle(Command{Action: action, Arg: arg, Result: resultCh})

	select {
	case res := <-resultCh:
		if res.OK {
			fmt.Fprintln(conn, strings.TrimSpace("OK "+res.Msg))
		} else {
			fmt.Fprintln(conn, strings.TrimSpace("ERR "+res.Msg))
		}
	case <-time.After(3 * time.Second):
		fmt.Fprintln(conn, "ERR timed out waiting for the running instance to respond")
	}
}

// SendCommand is the client side: connects to socketPath, sends action
// (plus arg, if any), and waits for the response. A connection failure
// means no peepoRun instance is currently listening there (not running,
// or the socket path is wrong) - the caller should treat that as "not
// running" rather than a generic error.
func SendCommand(socketPath, action, arg string) (ok bool, msg string, err error) {
	conn, err := net.DialTimeout("unix", socketPath, 1*time.Second)
	if err != nil {
		return false, "", err
	}
	defer conn.Close()

	line := action
	if arg != "" {
		line += " " + arg
	}
	if _, err := fmt.Fprintln(conn, line); err != nil {
		return false, "", err
	}

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return false, "", fmt.Errorf("no response from running instance")
	}
	resp := scanner.Text()
	switch {
	case strings.HasPrefix(resp, "OK"):
		return true, strings.TrimSpace(strings.TrimPrefix(resp, "OK")), nil
	case strings.HasPrefix(resp, "ERR"):
		return false, strings.TrimSpace(strings.TrimPrefix(resp, "ERR")), nil
	default:
		return false, "", fmt.Errorf("unexpected response: %q", resp)
	}
}
