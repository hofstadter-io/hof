package cmd

import (
	"fmt"
	"os"

	"log"
	"runtime/pprof"

	"strings"

	"github.com/hofstadter-io/hof/script/runtime"
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/config"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/ga"

)

var hofLong = `The High Code Framework`

func init() {

	RootCmd.PersistentFlags().StringSliceVarP(&(flags.RootPflags.Labels), "label", "l", nil, "Labels for use across all commands")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Config), "config", "", "", "Path to a hof configuration file")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Secret), "secret", "", "", "The path to a hof secret file")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.ContextFile), "context-file", "", "", "The path to a hof context file")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Context), "context", "", "", "The of an entry in the context file")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Global), "global", "", false, "Operate using only the global config/secret context")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Local), "local", "", false, "Operate using only the local config/secret context")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Package), "package", "p", "", "the Cue package context to use during this hof execution")
	RootCmd.PersistentFlags().IntVarP(&(flags.RootPflags.Verbose), "verbose", "v", 0, "set the verbosity of output")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Quiet), "quiet", "q", false, "turn off output and assume defaults at prompts")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Topic), "topic", "", "", "help topics for this command, 'list' will print available topics")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Example), "example", "", "", "examples for this command, 'list' will print available examples")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Tutorial), "tutorial", "", "", "tutorials for this command, 'list' will print available tutorials")
}

func RootPersistentPreRun(args []string) (err error) {

	config.Init()

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

		if flags.PrintSubject("Topics", "  ", flags.RootPflags.Topic, RootTopics) {
			return true
		}

		if flags.PrintSubject("Examples", "  ", flags.RootPflags.Example, RootExamples) {
			return true
		}

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

	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(DatamodelCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(TestCmd)
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(SecretCmd)
	RootCmd.AddCommand(ContextCmd)
	RootCmd.AddCommand(ReproduceCmd)
	RootCmd.AddCommand(JumpCmd)
	RootCmd.AddCommand(ReplCmd)
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
  init            β     create an empty workspace or initialize an existing directory to one
  gen             ✓     generate code, data, and config from your data models and designs
  mod             β     mod subcmd is a polyglot dependency management tool based on go mods
  test            α     test all sorts of things
  run             β     Hof Line Script (HLS) is a successor to bash and python based scripting
  config          β     manage local configurations
  secret          β     manage local secrets
  context         α     get, set, and use contexts

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  feedback        Ø     send feedback, bug reports, or any message

Learn more about hof and the _ you can do:
  each command has four flags, use 'list' as their arg
  to see available items on a command
    --help              print help message
    --topics            addtional help topics
    --examples          examples for the command
    --tutorials         tutorials for the command

(✓) command is generally available
(β) command is beta and ready for testing
(α) command is alpha and under developmenr
(Ø) command is null and yet to be implemented

Flags:
<<flag-usage>>
Use "hof [command] --help / -h" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.
`

var RootTopics = map[string]string{
	"main-topics": `There are several main topics:

hof gen  --topic list
hof mod  --topic list
hof test --topic list
hof run  --topic list
hof --topic script
`,
	"test-topic": `hello, this is a test topic.
please check out the others!
`,
}

var RootExamples = map[string]string{
	"test-example": `Thit is a test example, check out the subcommands, repo,  and website for more!
`,
}