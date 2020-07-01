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

var hofLong = `Polyglot Code Gereration Framework`

func init() {

	RootCmd.PersistentFlags().StringSliceVarP(&(flags.RootPflags.Labels), "label", "l", nil, "Labels for use across all commands")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Config), "config", "", "", "Path to a hof configuration file")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Secret), "secret", "", "", "The path to a hof secret file")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.ContextFile), "context-file", "", "", "The path to a hof context file")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Context), "context", "", "", "The of an entry in the context file")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Global), "global", "", false, "Operate using only the global config/secret context")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Local), "local", "", false, "Operate using only the local config/secret context")
	RootCmd.PersistentFlags().StringSliceVarP(&(flags.RootPflags.Input), "input", "i", nil, "input streams, depending on the command context")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.InputFormat), "input-format", "I", "", "input format, defaults to infered")
	RootCmd.PersistentFlags().StringSliceVarP(&(flags.RootPflags.Output), "output", "o", nil, "output streams, depending on the command context")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.OutputFormat), "output-format", "O", "", "output format, defaults to cue")
	RootCmd.PersistentFlags().StringSliceVarP(&(flags.RootPflags.Error), "error", "", nil, "error streams, depending on the command context")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.ErrorFormat), "error-format", "", "", "error format, defaults to cue")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Account), "account", "", "", "the account context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Billing), "billing", "", "", "the billing context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Project), "project", "", "", "the project context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Workspace), "workspace", "", "", "the workspace context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.DatamodelDir), "datamodel-dir", "", "", "directory for discovering resources")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.ResourcesDir), "resources-dir", "", "", "directory for discovering resources")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.RuntimesDir), "runtimes-dir", "", "", "directory for discovering runtimes")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Package), "package", "p", "", "the package context to use during this hof execution")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Errors), "all-errors", "E", false, "print all available errors")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Ignore), "ignore", "", false, "proceed in the presence of errors")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Simplify), "simplify", "S", false, "simplify output")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Trace), "trace", "", false, "trace cue computation")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Strict), "strict", "", false, "report errors for lossy mappings")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.Verbose), "verbose", "v", "", "set the verbosity of output")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.Quiet), "quiet", "q", false, "turn off output and assume defaults at prompts")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.ImpersonateAccount), "impersonate-account", "", "", "account to impersonate for this hof execution")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.TraceToken), "trace-token", "T", "", "used to help debug issues")
	RootCmd.PersistentFlags().StringVarP(&(flags.RootPflags.LogHTTP), "log-http", "", "", "used to help debug issues")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.RunWeb), "web", "", false, "run the command from the web ui")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.RunTUI), "tui", "", false, "run the command from the terminal ui")
	RootCmd.PersistentFlags().BoolVarP(&(flags.RootPflags.RunREPL), "repl", "", false, "run the command from the hof repl")
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

	Short: "Polyglot Code Gereration Framework",

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
	RootCmd.AddCommand(CloneCmd)
	RootCmd.AddCommand(DatamodelCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(RuntimesCmd)
	RootCmd.AddCommand(TestCmd)
	RootCmd.AddCommand(LabelCmd)
	RootCmd.AddCommand(LabelsetCmd)
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(AddCmd)
	RootCmd.AddCommand(CmdCmd)
	RootCmd.AddCommand(InfoCmd)
	RootCmd.AddCommand(CreateCmd)
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(SetCmd)
	RootCmd.AddCommand(EditCmd)
	RootCmd.AddCommand(DeleteCmd)
	RootCmd.AddCommand(DefCmd)
	RootCmd.AddCommand(EvalCmd)
	RootCmd.AddCommand(ExportCmd)
	RootCmd.AddCommand(FmtCmd)
	RootCmd.AddCommand(ImportCmd)
	RootCmd.AddCommand(TrimCmd)
	RootCmd.AddCommand(VetCmd)
	RootCmd.AddCommand(StCmd)
	RootCmd.AddCommand(AuthCmd)
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(SecretCmd)
	RootCmd.AddCommand(ContextCmd)
	RootCmd.AddCommand(StatusCmd)
	RootCmd.AddCommand(LogCmd)
	RootCmd.AddCommand(DiffCmd)
	RootCmd.AddCommand(BisectCmd)
	RootCmd.AddCommand(IncludeCmd)
	RootCmd.AddCommand(BranchCmd)
	RootCmd.AddCommand(CheckoutCmd)
	RootCmd.AddCommand(CommitCmd)
	RootCmd.AddCommand(MergeCmd)
	RootCmd.AddCommand(RebaseCmd)
	RootCmd.AddCommand(ResetCmd)
	RootCmd.AddCommand(TagCmd)
	RootCmd.AddCommand(FetchCmd)
	RootCmd.AddCommand(PullCmd)
	RootCmd.AddCommand(PushCmd)
	RootCmd.AddCommand(ProposeCmd)
	RootCmd.AddCommand(PublishCmd)
	RootCmd.AddCommand(RemotesCmd)
	RootCmd.AddCommand(ReproduceCmd)
	RootCmd.AddCommand(JumpCmd)
	RootCmd.AddCommand(UiCmd)
	RootCmd.AddCommand(TuiCmd)
	RootCmd.AddCommand(ShellCmd)
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

const RootCustomHelp = `hof - a polyglot tool for building software

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]


Initialize and create new hof workspaces:
  init            β     create an empty workspace or initialize an existing directory to one
  clone           β     clone a workspace or repository into a new directory

Model your designs, generate implementation, run or test anything:
  datamodel       α     create, view, diff, calculate / migrate, and manage your data models
  gen             ✓     generate code, data, and config from your data models and designs
  run             β     Hof Line Script (HLS) is a successor to bash and python based scripting
  test            α     test all sorts of things

Labels are used _ for _ (see also 'hof topic labels'):
  label           α     manage labels for resources and more
  labelset        α     group resources, datamodels, labelsets, and more

Learn more about hof and the _ you can do:
  each command has four flags, use 'list' as their arg
	to see available items on a command
		--help              print help message
		--topics            addtional help topics
		--examples          examples for the command
		--tutorials         tutorials for the command

Download modules, add instances or content, and manage runtimes:
  mod             β     mod subcmd is a polyglot dependency management tool based on go mods
  add             α     add dependencies and new components to the current module or workspace
  runtimes        α     work with runtimes (go, js, py, bash, docker, cloud-vms, k8s, custom)

Manage resources (see also 'hof topic resources'):
  info            α     print information about known resources
  create          α     create resources
  get             α     find and display resources
  set             α     find and configure resources
  edit            α     edit resources
  delete          α     delete resources

Configure, Unify, Execute (see also https://cuelang.org):
  cmd             α     run commands from the scripting layer and your _tool.cue files
  def             α     print consolidated definitions
  eval            α     print consolidated definitions
  export          α     export your data model to various formats
  fmt             α     formats code and files
  import          α     convert other formats and systems to hofland
  trim            α     cleanup code, configuration, and more
  vet             α     validate data
  st              α     recursive diff, merge, mask, pick, and query helpers for Cue

Manage logins, config, secrets, and context:
  auth            Ø     authentication subcommands
  config          β     manage local configurations
  secret          β     manage local secrets
  context         α     get, set, and use contexts

Examine workpsace history and state:
  status          α     show workspace information and status
  log             α     show workspace logs and history
  diff            α     show the difference between workspace versions
  bisect          α     use binary search to find the commit that introduced a bug

Grow, mark, and tweak your shared history (see also 'hof topic changesets'):
  include         α     include changes into the changeset
  branch          α     list, create, or delete branches
  checkout        α     switch branches or restore working tree files
  commit          α     record changes to the repository
  merge           α     join two or more development histories together
  rebase          α     reapply commits on top of another base tip
  reset           α     reset current HEAD to the specified state
  tag             α     create, list, delete or verify a tag object signed with GPG

Collaborate (see also 'hof topic collaborate'):
  fetch           α     download objects and refs from another repository
  pull            α     fetch from and integrate with another repository or a local branch
  push            α     update remote refs along with associated objects
  propose         α     propose to incorporate your changeset in a repository
  publish         α     publish a tagged version to a repository
  remotes         α     manage remote repositories

Local development commands:
  reproduce       Ø     Record, share, and replay reproducible environments and processes
  jump            α     Jumps help you do things with fewer keystrokes.
  ui              Ø     Run hof's local web ui
  tui             Ø     Run hof's terminal ui
  shell           α     Run hof's shell powered by HLS
  pprof                 go pprof by setting HOF_CPU_PROFILE="hof-cpu.prof" hof <cmd>


Send us feedback or say hello:
  feedback        Ø     send feedback, bug reports, or any message :]
                        you can also chat with us on https://gitter.im/hofstadter-io

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal

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
