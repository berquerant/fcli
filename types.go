package fcli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
)

var (
	ErrBadTargetFunction = fmt.Errorf("bad target function")
	ErrCallFailure       = fmt.Errorf("call failure")
)

// TargetFunction specifies a function for CLI subcommand.
type TargetFunction interface {
	Name() string
	// Call calls the function by flag arguments.
	Call(arguments []string) error
}

type targetFunction struct {
	f       any
	flags   []Flag
	flagSet *flag.FlagSet
	config  *Config
}

// NewTargetFunction makes a function able to be invoked by string slice arguments.
// F can be the function which is not variadic, no output parameters, not literal, not method
// and can have input parameters below:
//
//   int, int8, int16, int32, int64
//   uint, uint8, uint16, uint32, uint64
//   bool, string, float32, float64
//
// and the type which implements CustomFlagUnmarshaller.
// Default value is available if the type implements CustomFlagZeroer.
// Note: if pass the struct, pass as a pointer.
func NewTargetFunction(f any, opt ...Option) (TargetFunction, error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("%w not a function %v", ErrBadTargetFunction, f)
	}
	// find function
	fname, _ := GetFuncName(f)
	wrapErr := func(format string, v ...interface{}) error {
		return fmt.Errorf("%w %s %s", ErrBadTargetFunction, fname.FullName(), fmt.Sprintf(format, v...))
	}
	if t.IsVariadic() {
		return nil, wrapErr("variadic")
	}
	if t.NumOut() != 0 {
		return nil, wrapErr("has output parameters")
	}
	// read function AST
	fileLines, err := NewFileLines(fname.File())
	if err != nil {
		return nil, wrapErr("read file %w", err)
	}
	declSrc, err := NewFuncDeclCutter(fileLines, fname.Line()).CutFuncDecl()
	if err != nil {
		return nil, wrapErr("decl cut %w", err)
	}
	funcInfo, err := NewFuncInfo(declSrc)
	if err != nil {
		return nil, wrapErr("func info %w", err)
	}
	// generate flags from function
	flags := make([]Flag, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		p := t.In(i)
		ff, found := NewFlagFactory(p)
		if !found {
			return nil, wrapErr("unsupported parameter type %v", p)
		}
		flags[i] = ff(funcInfo.In(i).Name())
	}
	// apply options
	config := Config{
		ErrorHandling: flag.ExitOnError,
	}
	config.Apply(opt...)
	// init flags
	flagSet := flag.NewFlagSet(fname.String(), config.ErrorHandling)
	if doc := funcInfo.Doc(); doc != "" {
		flagSet.Usage = func() {
			fmt.Fprintf(os.Stderr, doc)
		}
	}
	for _, f := range flags {
		f.AddFlag(flagSet)
	}

	return &targetFunction{
		f:       f,
		flags:   flags,
		config:  &config,
		flagSet: flagSet,
	}, nil
}

func (s *targetFunction) Name() string { return s.flagSet.Name() }

func (s *targetFunction) Call(arguments []string) (rerr error) {
	defer func() {
		if err := recover(); err != nil {
			rerr = fmt.Errorf("%w recover %s %v", ErrCallFailure, s.flagSet.Name(), err)
		}
	}()

	if err := s.flagSet.Parse(arguments); err != nil {
		return fmt.Errorf("%w err %v", ErrCallFailure, err)
	}

	inputValues := make([]reflect.Value, len(s.flags))
	for i, f := range s.flags {
		v, err := f.ReflectValue()
		if err != nil {
			return fmt.Errorf("%w unwrap error %d th arg %s %v", ErrCallFailure, i+1, f.Name(), err)
		}
		inputValues[i] = v
	}
	_ = reflect.ValueOf(s.f).Call(inputValues)
	return nil
}