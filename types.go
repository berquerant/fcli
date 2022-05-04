package fcli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/berquerant/fcli/internal/ierrors"
)

var (
	ErrBadTargetFunction = errors.New("bad target function")
	ErrCallFailure       = errors.New("call failure")
)

// TargetFunction specifies a function for CLI subcommand.
type TargetFunction interface {
	Name() string
	// Call calls the function by flag arguments.
	// Returns ErrCallFailure if failed to call the function.
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
	fname, err := GetFuncName(f)
	if err != nil {
		return nil, fmt.Errorf("%w cannot get function name from %v", ErrBadTargetFunction, f)
	}

	wrapErr := ierrors.NewWrapperBuilder().
		Err(ErrBadTargetFunction).
		Msg("%s", fname.FullName()).
		Build()
	if t.IsVariadic() {
		return nil, wrapErr("variadic")
	}
	if t.NumOut() != 0 {
		return nil, wrapErr("has output parameters")
	}
	// read function AST
	funcInfo, err := BuildFuncInfo(fname.File(), fname.Line())
	if err != nil {
		return nil, wrapErr("build func info %w", err)
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
	config := NewConfigBuilder().
		ErrorHandling(flag.ExitOnError).
		CommandName(fname.String()).
		Build()
	config.Apply(opt...)
	// init flags
	flagSet := flag.NewFlagSet(
		config.CommandName.Get(),
		config.ErrorHandling.Get(),
	)
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
		config:  config,
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
