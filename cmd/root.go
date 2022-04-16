package cmd

import (
	"fmt"

	"github.com/Tea-n-Tech/chia-sync-helper/chia"
	"github.com/spf13/cobra"
)

const (
	cliHeader = ("        ___ _    _          ___                  _  _     _               \n" +
		"       / __| |_ (_)__ _ ___/ __|_  _ _ _  __ ___| || |___| |_ __  ___ _ _ \n" +
		"      | (__| ' \\| / _` |___\\__ \\ || | ' \\/ _|___| __ / -_) | '_ \\/ -_) '_|\n" +
		"       \\___|_||_|_\\__,_|   |___/\\_, |_||_\\__|   |_||_\\___|_| .__/\\___|_|  \n" +
		"                                 |__/                       |_|            \n")
	cliDescription = ("Chia sync helper identifies full node connections which are\n" +
		"behind in sync height and thus prevent us from syncing ourselves.")
	defaultHeightToleranceInBlocks = 5000
	defaultRunEveryMinutes         = 0
)

type CliArgs struct {
	HeightTolerance int64
	RunEveryMins    int64
}

var (
	cliArgs = CliArgs{}
)

var RootCmd = &cobra.Command{
	Use:   "chia-sync-helper",
	Short: "Chia sync helper removes connections not assisting with syncing.",
	Long:  cliHeader + "\n" + cliDescription,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf(cliHeader + "\n")
	},
	Run: func(cmd *cobra.Command, args []string) {
		chia.RunFullNodeCheck(cliArgs.RunEveryMins, cliArgs.HeightTolerance)
	},
}

func init() {
	flags := RootCmd.Flags()
	cliArgs.HeightTolerance = *flags.Int64P("height-tolerance",
		"t",
		defaultHeightToleranceInBlocks,
		("Every node whose height is lower than the current nodes height minus\n" +
			"'heigh-tolerance' will be removed."))
	cliArgs.RunEveryMins = *flags.Int64P("runs-every",
		"r",
		defaultRunEveryMinutes,
		"Runs the check indefinitely at the specified time interval in minutes.")
}
