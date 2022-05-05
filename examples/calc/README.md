# Example: calc

```
❯ ./calc -h
Error: command not found -h
Usage: calc {mult,pow,sum}
exit status 1

❯ ./calc pow -h
Usage of pow:
  -base int

  -exp int


❯ ./calc sum -h
sum prints the sum of args.
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

❯ ./calc pow -base 2 -exp 10
1024

```
