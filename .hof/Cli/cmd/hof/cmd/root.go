package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	// "github.com/spf13/viper"

	"github.com/hofstadter-io/hof/lib/runtime"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/cmd/hof/pflags"
)

var hofLong = `Polyglot Code Gereration Framework`

func init() {

	RootCmd.PersistentFlags().StringSliceVarP(&pflags.RootLabelsPflag, "label", "l", nil, "Labels for use across all commands")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootConfigPflag, "config", "", "", "Path to a hof configuration file")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootSecretPflag, "secret", "", "", "The path to a hof secret file")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootContextFilePflag, "context-file", "", "", "The path to a hof context file")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootContextPflag, "context", "", "", "The of an entry in the context file")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootGlobalPflag, "global", "", false, "Operate using only the global config/secret context")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootLocalPflag, "local", "", false, "Operate using only the local config/secret context")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootAccountPflag, "account", "", "", "the account context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootBillingPflag, "billing", "", "", "the billing context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootProjectPflag, "project", "", "", "the project context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootWorkspacePflag, "workspace", "", "", "the workspace context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootPackagePflag, "package", "p", "", "the package context to use during this hof execution")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootErrorsPflag, "all-errors", "E", false, "print all available errors")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootIgnorePflag, "ignore", "", false, "proceed in the presence of errors")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootSimplifyPflag, "simplify", "S", false, "simplify output")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootTracePflag, "trace", "", false, "trace cue computation")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootStrictPflag, "strict", "", false, "report errors for lossy mappings")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootVerbosePflag, "verbose", "v", "", "set the verbosity of output")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootQuietPflag, "quiet", "q", false, "turn off output and assume defaults at prompts")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootImpersonateAccountPflag, "impersonate-account", "", "", "account to impersonate for this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootTraceTokenPflag, "trace-token", "T", "", "used to help debug issues")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootLogHTTPPflag, "log-http", "", "", "used to help debug issues")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootRunUIPflag, "ui", "", false, "run the command from the web ui")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootRunTUIPflag, "tui", "", false, "run the command from the terminal ui")
}

func RootPersistentPreRun(args []string) (err error) {

	runtime.Init()

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

		ga.SendGaEvent("root", "<omit>", 0)

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

	help := func(cmd *cobra.Command, args []string) {
		fu := RootCmd.Flags().FlagUsages()
		rh := strings.Replace(RootCustomHelp, "<<flag-usage>>", fu, 1)
		fmt.Println(rh)
		fmt.Println(cmd.Name(), "hof", args)
	}
	usage := func(cmd *cobra.Command) error {
		fu := RootCmd.Flags().FlagUsages()
		rh := strings.Replace(RootCustomHelp, "<<flag-usage>>", fu, 1)
		fmt.Println(rh)
		return fmt.Errorf("unknown HOF command")
	}

	thelp := func(cmd *cobra.Command, args []string) {
		if RootCmd.Name() == cmd.Name() {
			ga.SendGaEvent("root/help", "<omit>", 0)
		}
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		if RootCmd.Name() == cmd.Name() {
			ga.SendGaEvent("root/usage", "<omit>", 0)
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
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(ModelsetCmd)
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(AddCmd)
	RootCmd.AddCommand(CmdCmd)
	RootCmd.AddCommand(LabelCmd)
	RootCmd.AddCommand(CreateCmd)
	RootCmd.AddCommand(ApplyCmd)
	RootCmd.AddCommand(GetCmd)
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
	RootCmd.AddCommand(TopicCmd)
	RootCmd.AddCommand(FeedbackCmd)
	RootCmd.AddCommand(DocCmd)
	RootCmd.AddCommand(JumpCmd)
	RootCmd.AddCommand(BuildCmd)
	RootCmd.AddCommand(UiCmd)
	RootCmd.AddCommand(TuiCmd)
	RootCmd.AddCommand(ReplCmd)
	RootCmd.AddCommand(HackCmd)
	RootCmd.AddCommand(GebCmd)
	RootCmd.AddCommand(LogoCmd)

}

const RootCustomHelp = `hof - a polyglot tool for building software

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]


Create or clone workspaces and datasets:
  + init          Create an empty Workspace or initialize an existing directory to one
  + clone         Clone a Workspace into a new directory

Model your world and generate implementation:
  modelset      create, view, migrate, and understand your modelsets.
  gen           generate code, data, and config

Download modules, add content, and run commands:
  mod           mod subcmd is a polyglot dependency management tool based on go mods
  + add           add dependencies and new components to the current module or workspace
  + cmd           Run commands from the scripting layer

Manage resources (see also 'hof topic resources'):
  + label         manage resource labels
  + create        create resources
  + apply         apply resource configuration
  + get           find and display resources
  + delete        delete resources

Configure, Unify, Execute (see also https://cuelang.org):
  + def           print consolidated definitions
  + eval          print consolidated definitions
  + export        export your data model to various formats
  + fmt           formats code and files
  + import        convert other formats and systems to hofland
  + trim          cleanup code, configuration, and more
  + vet           validate data
  st            Structural diff, merge, mask, pick, and query helpers for Cue


Manage logins, config, secrets, and context:
  + auth          authentication subcommands
  config        Manage local configurations
  secret        Manage local secrets
  + context       Get, set, and use contexts

Examine workpsace history and state:
  + status        Show workspace information and status
  + log           Show workspace logs and history
  + diff          Show the difference between workspace versions
  + bisect        Use binary search to find the commit that introduced a bug

Grow, mark, and tweak your shared history (see also 'hof topic changesets'):
  + include       Include changes into the changeset
  + branch        List, create, or delete branches
  + checkout      Switch branches or restore working tree files
  + commit        Record changes to the repository
  + merge         Join two or more development histories together
  + rebase        Reapply commits on top of another base tip
  + reset         Reset current HEAD to the specified state
  + tag           Create, list, delete or verify a tag object signed with GPG

Colloaborate (see also 'hof topic collaborate'):
  + fetch         Download objects and refs from another repository
  + pull          Fetch from and integrate with another repository or a local branch
  + push          Update remote refs along with associated objects
  + propose       Propose to include your changeset in a remote repository
  + reproduce     Record, share, and replay reproducible environments and processes
 
Local development commands:
  + doc           Generate and view documentation.
  + jump          Jumps help you get things done faster.
  + build         Build assets for modules and generated output
  + reproduce     Record, share, and replay reproducible environments and processes
  + ui            Run hof's local web ui
  + tui           Run hof's terminal ui
  + repl          Run hof's local REPL
  pprof         Go pprof by setting HOF_CPU_PROFILE="hof-cpu.prof" hof <cmd>


Send us feedback or say hello:
  + feedback      send feedback, bug reports, or any message :]

Additional commands:
  help          Help about any command
  topic         Additional information for various subjects and concepts
  update        Check for new versions and run self-updates
  version       Print detailed version information
  completion    Generate completion helpers for your terminal

Additional topics:
  schema, codegen, modeling, mirgrations
  resources, labels, context, querying
  workflow, changesets, collaboration

(+) command is yet to be implemented

Flags:
<<flag-usage>>
Use "hof [command] --help" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.
`
