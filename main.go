package main

import (
	"fmt"
	"log"
	"os"

	"runtime/pprof"

	"github.com/hofstadter-io/hof/commands"
)

func main() {

	// TODO, turn this into a flag and run if enabled, in root command persistent-pre-run, the flag should be a filename so we know where to write it
	f, err := os.Create("hof-cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
