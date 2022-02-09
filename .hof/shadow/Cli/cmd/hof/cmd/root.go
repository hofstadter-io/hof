package cmd

import (
	"fmt"
	"os"

	"log"
	"runtime/pprof"

	"strings"

	"github.com/hofstadter-io/hof/script/runtime"
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

var hofLong = `The High Code Framework`

func init() {

	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Package), "package", "p", "", "the Cue package context to use during execution")
	RootCmd.PersistentFlags().IntVarP(&(flags.RootPflags.Verbose), "verbose", "v", 0, "set the verbosity of output")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Quiet), "quiet", "q", false, "turn off output and assume defaults at prompts")
}

func RootPersistentPreRun(args []string) (err error) {

	return err
}

func RootPersistentPostRun(args []string) (err error) {

	WaitPrintUpdateAvailable()

	return err
}

var RootCmd = &cobra.Command{

	Use: "hof",

	Short: "The High Code Framework",

	Long: hofLong,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RootPersistentPreRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},

	PreRun: func(cmd *cobra.Command, args []string) {

	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RootPersistentPostRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func RootInit() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	help := func(cmd *cobra.Command, args []string) {
		if extra(cmd) {
			return
		}
		fu := RootCmd.Flags().FlagUsages()
		rh := strings.Replace(RootCustomHelp, "<<flag-usage>>", fu, 1)
		fmt.Println(rh)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		fu := RootCmd.Flags().FlagUsages()
		rh := strings.Replace(RootCustomHelp, "<<flag-usage>>", fu, 1)
		fmt.Println(rh)
		return fmt.Errorf("unknown hof command")
	}

	RootCmd.SetHelpFunc(help)
	RootCmd.SetUsageFunc(usage)

	RootCmd.AddCommand(UpdateCmd)

	RootCmd.AddCommand(VersionCmd)

	RootCmd.AddCommand(CompletionCmd)

	RootCmd.AddCommand(DatamodelCmd)
	RootCmd.AddCommand(FlowCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(FeedbackCmd)
	RootCmd.AddCommand(HackCmd)
	RootCmd.AddCommand(GebCmd)
	RootCmd.AddCommand(LogoCmd)

}

func RunExit() {
	if err := RunErr(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RunInt() int {
	if err := RunErr(); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func RunErr() error {

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

	RootInit()
	return RootCmd.Execute()
}

func CallTS(ts *runtime.Script, args []string) error {
	RootCmd.SetArgs(args)

	err := RootCmd.Execute()
	ts.Check(err)

	return err
}

const RootCustomHelp = `hof - the high code framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  datamodel             create, view, diff, calculate / migrate, and manage your data models
  flow                  run file(s) through the hof/flow DAG engine
  gen                   generate code, data, and config from your data models and designs
  mod                   mod subcmd is a polyglot dependency management tool based on go mods

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  feedback              send feedback, bug reports, or any message

Flags:
<<flag-usage>>
Use "hof [command] --help / -h" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.
`
