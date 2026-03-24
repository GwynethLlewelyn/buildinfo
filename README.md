# buildinfo
Retrieve common variables embedded in a Go executable.

That's all it does!

```go
package main

import (
	"fmt"

	"github.com/GwynethLlewelyn/buildinfo"
)

func main() {
	fmt.Println(buildinfo.String())
}
```

[![Go](https://github.com/GwynethLlewelyn/buildinfo/actions/workflows/go.yml/badge.svg)](https://github.com/GwynethLlewelyn/buildinfo/actions/workflows/go.yml)