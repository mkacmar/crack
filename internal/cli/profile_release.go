//go:build !debug

package cli

import "flag"

type profileConfig struct{}

func registerProfileFlags(_ *flag.FlagSet, _ *profileConfig) {}

func startProfiling(_ *profileConfig) (func(), error) { return nil, nil }

func profileUsage() string { return "" }
