# fcli

Package fcli provides utilities for function-based command-line tools.

## Usage

``` go
package main

import (
	"fmt"

	"github.com/berquerant/fcli"
)

func greet(name string) {
	fmt.Println("Hello,", name)
}

// bye prints Bye!
func bye() {
	fmt.Println("Bye!")
}

func main() {
	cli := fcli.NewCLI("do")
	_ = cli.Add(greet)
	_ = cli.Add(bye)
	_ = cli.Start()
}
```

```
❯ ./do
not enough arguments
Usage: do {bye,greet}

❯ ./do greet -h
Usage of greet:
  -name string

❯ ./do greet -name world
Hello, world

❯ ./do bye -h
bye prints Bye!

❯ ./do bye
Bye!
```

## Examples

[calc](./example/README.md)
