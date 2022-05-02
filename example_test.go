package fcli_test

import (
	"flag"
	"fmt"

	"github.com/berquerant/fcli"
)

func ExampleNewFlagFactory_basicTypes() {
	var (
		boolFlagF, _   = fcli.NewFlagFactory(false)
		stringFlagF, _ = fcli.NewFlagFactory("")

		boolFlag   = boolFlagF("b")
		stringFlag = stringFlagF("s")

		flagSet = flag.NewFlagSet("ff", flag.ContinueOnError)
	)

	boolFlag.AddFlag(flagSet)
	stringFlag.AddFlag(flagSet)

	if err := flagSet.Parse([]string{"-b", "-s", "basic"}); err != nil {
		panic(err)
	}
	b, _ := boolFlag.Unwrap()
	s, _ := stringFlag.Unwrap()
	fmt.Println(b)
	fmt.Println(s)
	// Output:
	// true
	// basic
}

func ExampleGetFuncName() {
	v, err := fcli.GetFuncName(fmt.Println)
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
	// Output:
	// Println
}

func ExampleNewFuncInfo() {
	const src = `// Hello greets you.
func Hello(name string) {
  println("Hello!", name)
}`
	info, err := fcli.NewFuncInfo(src)
	if err != nil {
		panic(err)
	}
	fmt.Println(info.Name())
	fmt.Println(info.NumIn())
	fmt.Println(info.In(0).Name())
	fmt.Println(info.Doc())
	// Output:
	// Hello
	// 1
	// name
	// Hello greets you.
}
