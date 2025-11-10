package cmd

import (
	"go-kit/tools/goctl/compare/testdata"
	"go-kit/tools/goctl/util/console"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare the goctl commands generated results between urfave and cobra",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		testdata.MustRun(dir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		console.Error("%+v", err)
	}
}
