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

	RootCmd.PersistentFlags().StringVarP(&pflags.RootConfigPflag, "config", "C", "", "Path to a hof configuration file")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootContextPflag, "context", "X", "", "The path to a hof creds file")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootAccountPflag, "account", "A", "", "the account context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootBillingPflag, "billing", "B", "", "the billing context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootProjectPflag, "project", "P", "", "the project context to use during this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootPackagePflag, "package", "p", "", "the package context to use during this hof execution")
	RootCmd.PersistentFlags().StringSliceVarP(&pflags.RootLabelsPflag, "label", "L", nil, "Labels for use across all commands")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootErrorsPflag, "all-errors", "E", false, "print all available errors")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootIgnorePflag, "ignore", "i", false, "proceed in the presence of errors")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootSimplifyPflag, "simplify", "s", false, "simplify output")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootTracePflag, "trace", "", false, "trace computation")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootStrictPflag, "strict", "", false, "report errors for lossy mappings")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootVerbosePflag, "verbose", "v", "", "set the verbosity of output")
	RootCmd.PersistentFlags().BoolVarP(&pflags.RootQuietPflag, "quiet", "q", false, "turn off output and assume defaults at prompts")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootImpersonateAccountPflag, "impersonate-account", "I", "", "account to impersonate for this hof execution")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootTraceTokenPflag, "trace-token", "", "", "used to help debug issues")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootLogHTTPPflag, "log-http", "", "", "used to help debug issues")
	RootCmd.PersistentFlags().StringVarP(&pflags.RootCredsPflag, "creds", "", "", "The path to a hof creds file")
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

func init() {

	help := func(cmd *cobra.Command, args []string) {
		fu := RootCmd.Flags().FlagUsages()
		rh := strings.Replace(RootCustomHelp, "<<flag-usage>>", fu, 1)
		fmt.Println(rh)
		fmt.Println("hof", args)
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
			ga.SendGaEvent("root/help", "<omit>", 0)
		}
		return usage(cmd)
	}
	RootCmd.SetHelpFunc(thelp)
	RootCmd.SetUsageFunc(tusage)

	// cobra.OnInitialize(initConfig)

	RootCmd.AddCommand(AuthCmd)
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(SecretCmd)
	RootCmd.AddCommand(CloneCmd)
	RootCmd.AddCommand(InitCmd)
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
	RootCmd.AddCommand(ModCmd)
	RootCmd.AddCommand(GenCmd)
	RootCmd.AddCommand(ModelCmd)
	RootCmd.AddCommand(AddCmd)
	RootCmd.AddCommand(CmdCmd)
	RootCmd.AddCommand(DefCmd)
	RootCmd.AddCommand(EvalCmd)
	RootCmd.AddCommand(ExportCmd)
	RootCmd.AddCommand(FmtCmd)
	RootCmd.AddCommand(ImportCmd)
	RootCmd.AddCommand(TrimCmd)
	RootCmd.AddCommand(VetCmd)
	RootCmd.AddCommand(LabelCmd)
	RootCmd.AddCommand(CreateCmd)
	RootCmd.AddCommand(ApplyCmd)
	RootCmd.AddCommand(GetCmd)
	RootCmd.AddCommand(DeleteCmd)
	RootCmd.AddCommand(TopicCmd)
	RootCmd.AddCommand(DevCmd)
	RootCmd.AddCommand(UiCmd)
	RootCmd.AddCommand(ReplCmd)
	RootCmd.AddCommand(HackCmd)

}

const RootCustomHelp = `hof - a polyglot tool for building software

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]


Create or clone workspaces and datasets:
  clone         Clone a Workspace into a new directory
  init          Create an empty Workspace or initialize an existing directory to one

Model your world and generate implementation:
  model         create, view, migrate, and understand your data models.
  gen           generate code, data, and config

Download modules, add content, and run commands:
  mod           mod subcmd is a polyglot dependency management tool based on go mods
  add           add dependencies and new components to the current module or workspace
  cmd           Run commands from the scripting layer

Configure, Unify, Execute (see also https://cuelang.org):
  def           print consolidated definitions
  eval          print consolidated definitions
  export        export your data model to various formats
  fmt           formats code and files
  import        convert other formats and systems to hofland
  trim          cleanup code, configuration, and more
  vet           validate data

Manage resources (see also 'hof topic resources'):
  label         manage resource labels
  create        create resources
  apply         apply resource configuration
  get           find and display resources
  delete        delete resources

Manage logins, config, and secrets:
  auth          authentication subcommands
  config        configuration subcommands
  secret        secret subcommands

Examine workpsace history and state:
  status        Show workspace information and status
  log           Show workspace logs and history
  diff          Show the difference between workspace versions
  bisect        Use binary search to find the commit that introduced a bug

Grow, mark, and tweak your shared history (see also 'hof topic changesets'):
  include       Include changes into the changeset
  branch        List, create, or delete branches
  checkout      Switch branches or restore working tree files
  commit        Record changes to the repository
  merge         Join two or more development histories together
  rebase        Reapply commits on top of another base tip
  reset         Reset current HEAD to the specified state
  tag           Create, list, delete or verify a tag object signed with GPG

Colloaborate (see also 'hof topic collaborate'):
  fetch         Download objects and refs from another repository
  pull          Fetch from and integrate with another repository or a local branch
  push           Update remote refs along with associated objects
  propose       Propose to include your changeset in a remote repository

Local development commands:
  dev           run hof's local dev setup
  ui            run hof's local web ui
  repl          run hof's REPL system

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

Flags:
<<flag-usage>>
Use "hof [command] --help" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.
`
