// Package app implements app-related functionality.
package app

import (
	"fmt"
	"net"

	"github.com/azazeal/fly/pkg/env"
)

// List returns the applications running in the current organization.
func List() ([]string, error) {
	return net.LookupTXT("_apps.internal")
}

// Regions reports the regions the current application is deployed to.
func Regions() ([]string, error) {
	return Deployments(env.Region())
}

// Deployments reports the regions the named application is deployed to.
func Deployments(app string) ([]string, error) {
	host := fmt.Sprintf("regions.%s.internal", app)

	return net.LookupTXT(host)
}
