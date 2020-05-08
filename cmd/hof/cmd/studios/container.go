package cmdstudios

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/container"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var containerLong = `Work with Hofstadter Studios containers`

func ContainerRun(args []string) (err error) {

	return err
}

var ContainerCmd = &cobra.Command{

	Use: "container",

	Aliases: []string{
		"cont",
	},

	Short: "Work with Hofstadter Studios containers",

	Long: containerLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ContainerRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	hf := ContainerCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		as := strings.Join(args, "/")
		ga.SendGaEvent(c+"/help", as, 0)
		hf(cmd, args)
	}
	ContainerCmd.SetHelpFunc(f)
	ContainerCmd.AddCommand(cmdcontainer.CallCmd)
	ContainerCmd.AddCommand(cmdcontainer.ListCmd)
	ContainerCmd.AddCommand(cmdcontainer.GetCmd)
	ContainerCmd.AddCommand(cmdcontainer.CreateCmd)
	ContainerCmd.AddCommand(cmdcontainer.UpdateCmd)
	ContainerCmd.AddCommand(cmdcontainer.DeployCmd)
	ContainerCmd.AddCommand(cmdcontainer.StatusCmd)
	ContainerCmd.AddCommand(cmdcontainer.PushCmd)
	ContainerCmd.AddCommand(cmdcontainer.PullCmd)
	ContainerCmd.AddCommand(cmdcontainer.ResetCmd)
	ContainerCmd.AddCommand(cmdcontainer.ShutdownCmd)
	ContainerCmd.AddCommand(cmdcontainer.DeleteCmd)
}
