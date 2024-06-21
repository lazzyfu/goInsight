package base

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
)

func read(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser, ch chan<- string) {
	reader := bufio.NewReader(std)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			ch <- readString
		}
	}
}

func Command(ctx context.Context, ch chan<- string, cmd string) error {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	// standard output
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	// standard error
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	// 一个需要读取stderr 另一个需要读取stdout
	wg.Add(2)
	go read(ctx, &wg, stderr, ch)
	go read(ctx, &wg, stdout, ch)
	// 这里一定要用start,而不是run
	if err := c.Start(); err != nil {
		return err
	}
	// 等待输出结束
	wg.Wait()
	// 获取退出状态
	var exitStatus int
	if err := c.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			exitStatus = ex.Sys().(syscall.WaitStatus).ExitStatus()
		}
	}
	if exitStatus != 0 {
		return fmt.Errorf("cmd exit status %d", exitStatus)
	}
	return nil
}
