package fcli_test

import (
	"flag"
	"fmt"
	"strings"

	"github.com/berquerant/fcli"
)

type targetArg struct {
	category string
	location string
}

func (*targetArg) UnmarshalFlag(v string) (fcli.CustomFlagUnmarshaller, error) {
	xs := strings.Split(v, ".")
	if len(xs) != 2 {
		return nil, fmt.Errorf("invalid format")
	}
	return &targetArg{
		category: xs[0],
		location: xs[1],
	}, nil
}

func queryCLI(key string, target *targetArg) {
	fmt.Println(key)
	fmt.Println(target.category)
	fmt.Println(target.location)
}

func ExampleTargetFunction() {
	f, err := fcli.NewTargetFunction(queryCLI, fcli.WithErrorHandling(flag.ContinueOnError))
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Name())
	if err := f.Call([]string{"-key", "fall", "-target", "human.mars"}); err != nil {
		panic(err)
	}
	// Output:
	// queryCLI
	// fall
	// human
	// mars
}
