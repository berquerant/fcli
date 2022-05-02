package fcli_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/berquerant/fcli"
)

type customJSONStruct struct {
	World  string `json:"world"`
	Number int    `json:"number"`
}

func (*customJSONStruct) UnmarshalFlag(v string) (fcli.CustomFlagUnmarshaller, error) {
	var x customJSONStruct
	if err := json.Unmarshal([]byte(v), &x); err != nil {
		return nil, err
	}
	return &x, nil
}

type customStringSlice []string

func (customStringSlice) UnmarshalFlag(v string) (fcli.CustomFlagUnmarshaller, error) {
	return customStringSlice(strings.Split(v, ",")), nil
}

type customStringZero string

func (customStringZero) UnmarshalFlag(v string) (fcli.CustomFlagUnmarshaller, error) {
	return customStringZero(v), nil
}

func (customStringZero) FlagZero() fcli.CustomFlagUnmarshaller {
	return customStringZero("ZERO")
}

func ExampleNewFlagFactory_customFlag() {
	var (
		jsonFlagF, _  = fcli.NewFlagFactory(new(customJSONStruct))
		sliceFlagF, _ = fcli.NewFlagFactory(customStringSlice(nil))
		zeroFlagF, _  = fcli.NewFlagFactory(customStringZero(""))

		jsonFlag  = jsonFlagF("j")
		sliceFlag = sliceFlagF("s")
		zeroFlag  = zeroFlagF("z")

		flagSet = flag.NewFlagSet("ff", flag.ContinueOnError)
	)

	jsonFlag.AddFlag(flagSet)
	sliceFlag.AddFlag(flagSet)
	zeroFlag.AddFlag(flagSet)

	if err := flagSet.Parse([]string{"-j", `{"world":"fog","number":1}`, "-s", "sig,light,back"}); err != nil {
		panic(err)
	}
	j, _ := jsonFlag.Unwrap()
	s, _ := sliceFlag.Unwrap()
	z, _ := zeroFlag.Unwrap()
	fmt.Println(j.(*customJSONStruct).World)
	fmt.Println(j.(*customJSONStruct).Number)
	for _, x := range s.(customStringSlice) {
		fmt.Println(x)
	}
	fmt.Println(z)
	// Output:
	// fog
	// 1
	// sig
	// light
	// back
	// ZERO
}
