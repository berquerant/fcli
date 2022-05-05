package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/berquerant/fcli"
)

func sqrt(x float64) error {
	if x < 0 {
		return errors.New("negative")
	}
	fmt.Println(math.Sqrt(x))
	return nil
}

func wait(ctx context.Context, durationMS int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Duration(durationMS) * time.Millisecond):
		return nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	cli := fcli.NewCLI("ctx")
	cli.OnError(func(err error) int {
		if errors.Is(err, fcli.ErrCLICommandNotFound) || errors.Is(err, fcli.ErrCLINotEnoughArguments) {
			return fcli.Cusage
		}
		return fcli.Cerror
	})
	_ = cli.Add(sqrt)
	_ = cli.Add(wait)
	if err := cli.StartWithContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
