package cmd

import (
	"errors"

	"github.com/dovydas1928/tb"
	helper "github.com/dovydas1928/tb/cmd"
	"github.com/spf13/cobra"
)

func getCommon() tb.Common {
	return tb.NewCommonRemote(getConnection())
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := getContext()
		defer cancel()

		commonLocal := tb.NewCommonLocal()
		localVer, localRev, _ := commonLocal.GetVersion(ctx)
		helper.PrintVersion("tb-client", localVer, localRev)

		commonRemote := getCommon()

		remoteVer, remoteRev, err := commonRemote.GetVersion(ctx)
		helper.ExitOnError("could not get server version", err)

		helper.PrintVersion("tb-server", remoteVer, remoteRev)
	},
}

var modprobeCmd = &cobra.Command{
	Use:   "modprobe [module]",
	Short: "Run modbprobe command on server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("command requires [module] argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := getContext()
		defer cancel()

		commonRemote := getCommon()
		err := commonRemote.Modprobe(ctx, args[0])
		helper.ExitOnError("unable to modprobe", err)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd, modprobeCmd)
}
