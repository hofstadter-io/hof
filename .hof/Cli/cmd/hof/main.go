package main

import (
	"fmt"
	"os"

	"log"
	"runtime/pprof"

	"github.com/hofstadter-io/hof/cmd/hof/cmd"
)

func main() {

	if fn := os.Getenv("HOF_CPU_PROFILE"); fn != "" {
		f, err := os.Create(fn)
		if err != nil {
			log.Fatal("Could not create file for CPU profile:", err)
		}
		defer f.Close()

		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal("Could not start CPU profile process:", err)
		}

		defer pprof.StopCPUProfile()
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
