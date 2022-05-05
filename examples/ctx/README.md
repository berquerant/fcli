# Example: ctx

Return an error.

```
❯ ./ctx sqrt -x -1
negative
exit status 1
```

Pass context.

```
❯ ./ctx wait -durationMS 300
context deadline exceeded
exit status 1
```

Error but exit status is 0, because `fcli.Cusage` returned by the function set by `cli.OnError()`.

```
❯ ./ctx
Error: not enough arguments
Usage: ctx {sqrt,wait}
```
