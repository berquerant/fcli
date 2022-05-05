// Package fcli provides utilities for function-based command-line tools.
package fcli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/berquerant/fcli/internal/logger"
)

var (
	ErrCLINotEnoughArguments = errors.New("not enough arguments")
	ErrCLICommandNotFound    = errors.New("command not found")

	// NilUsage is noop.
	// Disable Usage of CLI by CLI.Usage(NilUsage).
	NilUsage       = func() {}
	DefaultOnError = func(_ error) int { return Cusage | Cerror }
)

const (
	Cusage = 1 << iota // print usage
	Cerror             // return error
)

// CLI is the function-based subcommands set.
type CLI interface {
	// Start parses arguments and calls proper function.
	// if arguments is nil, reads os.Args.
	Start(arguments ...string) error
	StartWithContext(ctx context.Context, arguments ...string) error
	// Add adds a subcommand.
	// See NewTargetFunction.
	Add(f any, opt ...Option) error
	// Usage sets a function to print usage.
	Usage(func())
	// OnError sets a function called when command function returned an error.
	// If onError return Cusage then print usage.
	// If onError return Cerror then return the error.
	OnError(onError func(error) int)
}

func NewCLI(name string, opt ...Option) CLI {
	return &cliMap{
		name:     name,
		commands: map[string]TargetFunction{},
		onError:  DefaultOnError,
	}
}

type cliMap struct {
	name     string
	usage    func()
	onError  func(error) int
	commands map[string]TargetFunction
}

func (s *cliMap) StartWithContext(ctx context.Context, arguments ...string) error {
	if err := s.start(ctx, arguments...); err != nil {
		r := s.onError(err)
		if r&Cusage != 0 {
			s.printUsage(err)
		}
		if r&Cerror != 0 {
			return err
		}
	}
	return nil
}

func (s *cliMap) Start(arguments ...string) error {
	return s.StartWithContext(context.Background(), arguments...)
}

func (s *cliMap) start(ctx context.Context, arguments ...string) error {
	args, ok := func() ([]string, bool) {
		if len(arguments) == 0 {
			if len(os.Args) < 2 {
				return nil, false
			}
			return os.Args[1:], true
		}
		return arguments, true
	}()
	if !ok {
		return ErrCLINotEnoughArguments
	}

	cmd, ok := s.commands[args[0]]
	if !ok {
		return fmt.Errorf("%w %s", ErrCLICommandNotFound, args[0])
	}
	logger.Debug("Call %s with %#v", cmd.Name(), args[1:])
	return cmd.CallWithContext(ctx, args[1:])
}

func (s *cliMap) printUsage(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	if s.usage != nil {
		s.usage()
		return
	}
	s.defaultUsage()
}

func (s *cliMap) defaultUsage() {
	var (
		i  int
		ss = make([]string, len(s.commands))
	)
	for k := range s.commands {
		ss[i] = k
		i++
	}
	sort.Strings(ss)
	fmt.Fprintf(os.Stderr, "Usage: %s {%s}\n", s.name, strings.Join(ss, ","))
}

func (s *cliMap) Add(f any, opt ...Option) error {
	t, err := NewTargetFunction(f, opt...)
	if err != nil {
		return err
	}
	logger.Debug("Add command %s %#v to %s", t.Name(), f, s.name)
	s.commands[t.Name()] = t
	return nil
}

func (s *cliMap) Usage(usage func())              { s.usage = usage }
func (s *cliMap) OnError(onError func(error) int) { s.onError = onError }
