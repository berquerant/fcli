package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/berquerant/fcli"
)

type intList []int

func (intList) UnmarshalFlag(v string) (fcli.CustomFlagUnmarshaller, error) {
	ss := strings.Split(v, ",")
	xs := make([]int, len(ss))
	for i, s := range ss {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		xs[i] = v
	}
	return intList(xs), nil
}

func (intList) FlagZero() fcli.CustomFlagUnmarshaller {
	return intList([]int{})
}

// sum prints the sum of args.
func sum(args intList) {
	var s int
	for _, a := range args {
		s += a
	}
	fmt.Println(s)
}

type comp complex128

func (comp) UnmarshalFlag(v string) (fcli.CustomFlagUnmarshaller, error) {
	ss := strings.Split(v, ",")
	if len(ss) == 2 {
		a, err := strconv.Atoi(ss[0])
		if err != nil {
			return nil, err
		}
		b, err := strconv.Atoi(ss[1])
		if err != nil {
			return nil, err
		}
		return comp(complex(float64(a), float64(b))), nil
	}
	return nil, fmt.Errorf("invalid format")
}

// mult multiplies two complex numbers.
func mult(a, b comp) {
	c := a * b
	fmt.Println(c)
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cli := fcli.NewCLI("calc")
	fail(cli.Add(sum))
	fail(cli.Add(mult))
	_ = cli.Start()
}
