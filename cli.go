// Package fcli provides utilities for function-based command-line tools.
package fcli

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	ErrCLINotEnoughArguments = fmt.Errorf("not enough arguments")
	ErrCLICommandNotFound    = fmt.Errorf("command not found")

	// NilUsage is noop.
	// Disable Usage of CLI by CLI.Usage(NilUsage).
	NilUsage = func() {}
)

// CLI is the function-based subcommands set.
type CLI interface {
	// Start parses arguments and calls proper function.
	// if arguments is nil, reads os.Args.
	Start(arguments ...string) error
	// Add adds a subcommand.
	// See NewTargetFunction.
	Add(f any, opt ...Option) error
	// Usage sets a function to print usage.
	Usage(func())
}

func NewCLI(name string) CLI {
	return &cliMap{
		name:     name,
		commands: map[string]TargetFunction{},
	}
}

type cliMap struct {
	name     string
	usage    func()
	commands map[string]TargetFunction
}

func (s *cliMap) Start(arguments ...string) error {
	if err := s.start(arguments...); err != nil {
		s.printUsage(err)
		return err
	}
	return nil
}

func (s *cliMap) start(arguments ...string) error {
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
	return cmd.Call(args[1:])
}

func (s *cliMap) printUsage(err error) {
	fmt.Fprintln(os.Stderr, err)
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
	s.commands[t.Name()] = t
	return nil
}

func (s *cliMap) Usage(usage func()) { s.usage = usage }
