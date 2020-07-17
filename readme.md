# Repl.it Database Go client

[![PkgGoDev](https://pkg.go.dev/badge/github.com/replit/database-go)](https://pkg.go.dev/github.com/replit/database-go)

The easiest way to use Repl.it Database from your Go repls.
[Try it out on Repl.it!](https://repl.it/@kochman/Database-Go-example)

```
package main

import (
  "fmt"
  "github.com/replit/database-go"
)

func main() {
  database.Set("key", "value")
  val, _ := database.Get("key")
  fmt.Println(val)
}
```

[View the docs](https://pkg.go.dev/github.com/replit/database-go) for more info
about how to use the client to interact with Database.
