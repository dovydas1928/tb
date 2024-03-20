package cmd

import (
	"context"

	"github.com/dovydas1928/tb"
	helper "github.com/dovydas1928/tb/cmd"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		common := tb.NewCommonLocal()
		ver, rev, _ := common.GetVersion(context.Background())
		helper.PrintVersion("tb-server", ver, rev)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
