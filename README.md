[![build](https://github.com/azazeal/fly/actions/workflows/build.yml/badge.svg)](https://github.com/azazeal/fly/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/azazeal/fly.svg)](https://pkg.go.dev/github.com/azazeal/fly)

# fly

Package env implements functionality for when dealing with fly.io runtime environment.

## Usage

```go
package main

import (
	"fmt"

	"github.com/azazeal/fly/pkg/env"
)

func main() {
	const format = `running on fly: %t
---
$FLY_APP_NAME: %s,
$FLY_ALLOC_ID: %s,
$FLY_PUBLIC_IP: %s,
$FLY_REGION: %s,
`
	fmt.Printf(format,
		env.IsSet(),
		env.AppName(),
		env.AllocID(),
		env.PublicIP(),
		env.Region(),
	)
}
```
