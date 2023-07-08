package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/charmbracelet/log"
)

type Runner struct{}

type Client struct {
	logger *log.Logger
	ch     chan string
	ctx    context.Context
}

func NewClient(ctx context.Context, opts log.Options, ch chan string) *Client {
	return &Client{
		logger: log.NewWithOptions(os.Stdout, opts),
		ch:     ch,
		ctx:    ctx,
	}
}

func (c *Client) output() {
	ticker := time.NewTicker(time.Millisecond * 300)
	var count = 1
	for {
		select {
		case t := <-ticker.C:
			c.logger.Info("Client output")
			c.ch <- fmt.Sprintf("Client output: %v", t)
			count++
		case <-c.ctx.Done():
			c.logger.Info("Client exit")
			return
		}
	}
}

func (r *Runner) Run(ctx context.Context, opts log.Options, idx int) {
	logger := log.NewWithOptions(os.Stdout, opts)

	ch := make(chan string)
	for i := 0; i < 2; i++ {
		go (NewClient(ctx, opts, ch)).output()
	}

	for {
		select {
		case t := <-ch:
			logger.Debug("Receive", "idx", idx, "msg", t)
		case <-ctx.Done():
			logger.Info("Done")
			return
		}
	}
}

func TestOutput(t *testing.T) {
	opts := log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
		TimeFormat:      time.DateTime,
		Level:           log.DebugLevel,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go (&Runner{}).Run(ctx, opts, 1)

	stoper := make(chan os.Signal, 3)
	signal.Notify(stoper, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-stoper
	cancel()
	time.Sleep(time.Second)
	log.Info("程序已停止")
}

func TestMultiple(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "MAIN",
	})
	for i := 0; i < 2; i++ {
		go func(i int) {
			subLog := logger.WithPrefix(fmt.Sprintf("IDX:%d", i))
			_ = subLog
			ticker := time.NewTicker(time.Millisecond * 500)
			count := 1
			for {
				select {
				case <-ticker.C:
					subLog.Infof("Loop %d times", count)
				case <-ctx.Done():
					subLog.Warnf("Go routines: %d", count)
					return
				}
				count++
			}
		}(i)
	}

	logger.Info("exit...")
	<-ctx.Done()
}
