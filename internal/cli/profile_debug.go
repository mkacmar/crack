//go:build debug

package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

type profileConfig struct {
	cpuProfile string
	memProfile string
}

func registerProfileFlags(fs *flag.FlagSet, cfg *profileConfig) {
	fs.StringVar(&cfg.cpuProfile, "cpuprofile", "", "")
	fs.StringVar(&cfg.memProfile, "memprofile", "", "")
}

func startProfiling(cfg *profileConfig) (stop func(), err error) {
	var cpuFile *os.File

	if cfg.cpuProfile != "" {
		cpuFile, err = os.Create(cfg.cpuProfile)
		if err != nil {
			return nil, fmt.Errorf("could not create CPU profile: %w", err)
		}
		if err = pprof.StartCPUProfile(cpuFile); err != nil {
			cpuFile.Close()
			return nil, fmt.Errorf("could not start CPU profile: %w", err)
		}
	}

	stop = func() {
		if cpuFile != nil {
			pprof.StopCPUProfile()
			cpuFile.Close()
		}
		if cfg.memProfile != "" {
			f, err := os.Create(cfg.memProfile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: could not create memory profile: %v\n", err)
				return
			}
			defer f.Close()
			runtime.GC()
			if err := pprof.WriteHeapProfile(f); err != nil {
				fmt.Fprintf(os.Stderr, "Error: could not write memory profile: %v\n", err)
			}
		}
	}
	return stop, nil
}

func profileUsage() string {
	return `
Profiling options:
      --cpuprofile string     Write CPU profile to file
      --memprofile string     Write memory profile to file
`
}
