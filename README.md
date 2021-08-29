[![Build Status](https://github.com/azazeal/fly/actions/workflows/build.yml/badge.svg)](https://github.com/azazeal/fly/actions/workflows/build.yml)
[![Coverage Report](https://coveralls.io/repos/github/azazeal/fly/badge.svg?branch=master)](https://coveralls.io/github/azazeal/fly?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/azazeal/fly.svg)](https://pkg.go.dev/github.com/azazeal/fly)

# fly

Package fly and its sub-packages implement helpers for Go apps running on
[fly.io](https://fly.io).

## Usage

```go
package main

import (
	"context"
	"log"
	"strings"

	"github.com/azazeal/fly/dns"
	"github.com/azazeal/fly/env"
)

func main() {
	if !env.IsSet() {
		log.Fatal("not running on fly")
	}

	const format = `running on fly: %t
---
$FLY_APP_NAME: %s,
$FLY_ALLOC_ID: %s,
$FLY_PUBLIC_IP: %s,
$FLY_REGION: %s,
`
	log.Printf(format,
		env.IsSet(),
		env.AppName(),
		env.AllocID(),
		env.PublicIP(),
		env.Region(),
	)

	apps, err := dns.Apps(context.TODO())
	if err != nil {
		log.Fatalf("failed determining apps: %v", apps)
	}

	log.Printf("fly apps in our organization: %s", strings.Join(apps, ", "))
}
```
