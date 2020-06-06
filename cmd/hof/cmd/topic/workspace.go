package cmdtopic

import (
	"fmt"

	"github.com/spf13/cobra"

	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var workspaceLong = `help for learning about workspaces and workflows`

var WorkspaceCmd = &cobra.Command{

	Use: "workspace",

	Short: "help for learning about workspaces and workflows",

	Long: workspaceLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {

	help := func(cmd *cobra.Command, args []string) {
		fu := WorkspaceCmd.Flags().FlagUsages()
		ch := strings.Replace(WorkspaceCustomHelp, "<<flag-usage>>", fu, 1)
		fmt.Println(ch)

	}
	usage := func(cmd *cobra.Command) error {
		fu := WorkspaceCmd.Flags().FlagUsages()
		ch := strings.Replace(WorkspaceCustomHelp, "<<flag-usage>>", fu, 1)
		fmt.Println(ch)

		return fmt.Errorf("unknown command %q", cmd.Name())
	}

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	WorkspaceCmd.SetHelpFunc(thelp)
	WorkspaceCmd.SetUsageFunc(tusage)

}

const WorkspaceCustomHelp = `These are common Workspace and Dataset commands used in various situations:

start a working area (see also: git help tutorial)
	 clone      Clone a Workspace into a new directory
	 init       Create an empty Workspace or reinitialize an existing one

examine the history and state
	 status     Show the working tree status
	 log        Show commit logs
	 diff       Show changes between commits, commit and working tree, etc
	 bisect     Use binary search to find the commit that introduced a bug

grow, mark and tweak your shared history
	 add        Add file changes to the index
	 branch     List, create, or delete branches
	 checkout   Switch branches or restore working tree files
	 commit     Record changes to the repository
	 merge      Join two or more development histories together
	 rebase     Reapply commits on top of another base tip
	 reset      Reset current HEAD to the specified state
	 tag        Create, list, delete or verify a tag object signed with GPG

collaborate and work with remote members
	 fetch      Download objects and refs from another repository
	 pull       Fetch from and integrate with another repository or a local branch
	 push       Update remote refs along with associated objects
	 propose    Propose to include your changeset in a remote repository
`
