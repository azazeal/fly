// Package env implements functionality for when dealing with fly.io
// environment variables.
package env

import "os"

// The set of environment variable keys.
const (
	// AppNameKey denotes the name of the environment variable which reports
	// the application's name.
	AppNameKey = "FLY_APP_NAME"

	// AllocIDKey denotes the name of the environment variable which reports
	// the instance's ID.
	AllocIDKey = "FLY_ALLOC_ID"

	// PublicIPKey denotes the name of the environment variable which reports
	// the instance's public IP address.
	PublicIPKey = "FLY_PUBLIC_IP"

	// RegionKey denotes the name of the environment variable which reports
	// the instance's region.
	RegionKey = "FLY_REGION"
)

var (
	keys = []string{AppNameKey, AllocIDKey, PublicIPKey, RegionKey}

	lookups = []func() (string, bool){
		LookupAppName,
		LookupAllocID,
		LookupPublicIP,
		LookupRegion,
	}
)

// IsSet reports whether all fly-related environment variables are defined.
func IsSet() bool {
	for _, fn := range lookups {
		if _, ok := fn(); !ok {
			return false
		}
	}

	return true
}

// Map returns a map containing all the defined fly-related environment
// variables.
//
// In case no fly-related environement variable is set, the returned map
// will be nil.
func Map() (kv map[string]string) {
	for _, key := range keys {
		v, ok := lookup(key)
		if !ok {
			continue
		}

		if kv == nil {
			kv = make(map[string]string, len(keys))
		}
		kv[key] = v
	}

	return
}

var lookup = os.LookupEnv

func get(name string) string {
	v, _ := lookup(name)
	return v
}

// AppName is shorthand for os.Getenv(AppNameKey).
func AppName() string {
	return get(AppNameKey)
}

// LookupAppName is shorthand for os.LookupEnv(AppNameKey).
func LookupAppName() (string, bool) {
	return lookup(AppNameKey)
}

// AllocID is shorthand for os.Getenv(AllocIDKey).
func AllocID() string {
	return get(AllocIDKey)
}

// LookupAllocID is shorthand for os.LookupEnv(AllocIDKey).
func LookupAllocID() (string, bool) {
	return lookup(AllocIDKey)
}

// PublicIP is shorthand for os.Getenv(PublicIPKey).
func PublicIP() string {
	return get(PublicIPKey)
}

// LookupPublicIP is shorthand for os.LookupEnv(PublicIPKey).
func LookupPublicIP() (string, bool) {
	return lookup(PublicIPKey)
}

// Region is shorthand for os.Getenv(RegionKey).
func Region() (v string) {
	return get(RegionKey)
}

// LookupRegion is shorthand for os.LookupEnv(RegionKey).
func LookupRegion() (string, bool) {
	return lookup(RegionKey)
}
