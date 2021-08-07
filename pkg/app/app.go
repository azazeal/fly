// Package app implements app-related functionality.
package app

import (
	"context"
	"fmt"
	"net"

	"github.com/azazeal/fly/pkg/env"
)

var lookupTXT = net.DefaultResolver.LookupTXT

// List returns the applications running in the current organization.
func List(ctx context.Context) ([]string, error) {
	return lookupTXT(ctx, "_apps.internal")
}

// Regions reports the regions the current application is deployed to.
func Regions(ctx context.Context) ([]string, error) {
	return Deployments(ctx, env.AppName())
}

// Deployments reports the regions the named application is deployed to.
func Deployments(ctx context.Context, app string) ([]string, error) {
	host := fmt.Sprintf("regions.%s.internal", app)

	return lookupTXT(ctx, host)
}
