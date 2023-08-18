package cmd

import (
	"fmt"
	"os"

	"log"
	"runtime/pprof"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var hofLong = `The High Code Framework`

func init() {

	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Package), "package", "p", "", "the Cue package context to use during execution")
	RootCmd.PersistentFlags().StringSliceVarP(&(flags.RootPflags.Tags), "tags", "t", nil, "@tags() to be injected into CUE code")
	RootCmd.PersistentFlags().StringSliceVarP(&(flags.RootPflags.Path), "path", "l", nil, "CUE expression for single path component when placing data files")
	RootCmd.PersistentFlags().IntVarP(&(flags.RootPflags.Verbosity), "verbosity", "v", 0, "set the verbosity of output")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.IncludeData), "include-data", "", false, "auto include all data files found with cue files")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.InjectEnv), "inject-env", "", false, "inject all ENV VARs as default tag vars")
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

		ga.SendCommandPath("root")

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

	thelp := func(cmd *cobra.Command, args []string) {
		if RootCmd.Name() == cmd.Name() {
			ga.SendCommandPath("root help")
		}
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		if RootCmd.Name() == cmd.Name() {
			ga.SendCommandPath("root usage")
		}
		return usage(cmd)
	}
	RootCmd.SetHelpFunc(thelp)
	RootCmd.SetUsageFunc(tusage)

	RootCmd.AddCommand(UpdateCmd)

	RootCmd.AddCommand(VersionCmd)

	RootCmd.AddCommand(CompletionCmd)

	RootCmd.AddCommand(CreateCmd)
	RootCmd.AddCommand(DatamodelCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(FlowCmd)
	RootCmd.AddCommand(StCmd)
	RootCmd.AddCommand(FmtCmd)
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(DefCmd)
	RootCmd.AddCommand(EvalCmd)
	RootCmd.AddCommand(ExportCmd)
	RootCmd.AddCommand(VetCmd)
	RootCmd.AddCommand(ChatCmd)
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(FeedbackCmd)

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

const RootCustomHelp = `hof - the high code framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  chat                  co-create with AI (alpha)
  create                starter kits or blueprints from any git repo
  datamodel             manage, diff, and migrate your data models
  def                   print consolidated CUE definitions
  eval                  evaluate and print CUE configuration
  export                output data in a standard format
  flow                  run workflows and tasks powered by CUE
  fmt                   format any code and manage the formatters
  gen                   CUE powered code generation
  mod                   CUE module dependency management
  st                    apply CUE transformations in bulk
  vet                   validate data with CUE

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  feedback              open an issue or discussion on GitHub

Flags:
<<flag-usage>>
Use "hof [command] --help / -h" for more information about a command.
`
