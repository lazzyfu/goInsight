package base

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

func read(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser, ch chan<- string) {
	defer wg.Done()
	defer std.Close()

	reader := bufio.NewReader(std)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			ch <- line
		}
	}
}

func Command(ctx context.Context, ch chan<- string, binary string, args []string) error {
	c := exec.CommandContext(ctx, binary, args...)

	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go read(ctx, &wg, stdout, ch)
	go read(ctx, &wg, stderr, ch)

	if err := c.Start(); err != nil {
		return err
	}

	wg.Wait()

	if err := c.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			if ws, ok := ex.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("cmd exit status %d", ws.ExitStatus())
			}
		}
		return err
	}

	return nil
}

func RenderCommandForLog(binary string, args []string) string {
	logArgs := make([]string, len(args))
	copy(logArgs, args)
	for i, arg := range logArgs {
		if strings.HasPrefix(arg, "--password=") {
			logArgs[i] = "--password=..."
		}
	}
	return strings.TrimSpace(strings.Join(append([]string{binary}, logArgs...), " "))
}
