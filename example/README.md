# Example: calc

Usage.

```
❯ ./calc -h
command not found -h
Usage: calc {mult,sum}

❯ ./calc mult -h
mult multiplies two complex numbers.

❯ ./calc sum -h
sum prints the sum of args.
```

No arguments

```
❯ ./calc
not enough arguments
Usage: calc {mult,sum}
```

`sum` without arguments prints 0 because the default value of `intList` is `intList([]int{})` from `intList.FlagZero()`.

```
❯ ./calc sum
0
```

`mult` without arguments fails to run because the default value of `comp` is zero value.

```
❯ ./calc mult
call failure recover mult reflect: Call using zero Value argument
Usage: calc {mult,sum}
```

Normal cases:

```
❯ ./calc sum -args 1
1

❯ ./calc sum -args 1,2,3,4
10

❯ ./calc mult -a 1,2 -b 3,4
(-5+10i)
```
