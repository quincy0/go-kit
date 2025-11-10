package doc

import (
	"github.com/spf13/cobra"
)

// Cmd describes a upgrade command.
var (
	api string
	app string

	Cmd = &cobra.Command{
		Use:   "doc",
		Short: "push api docs to yapi",
		RunE:  Sync,
	}
)

func init() {
	Cmd.Flags().StringVar(&app, "app", "", "select app to push yapi")
	Cmd.Flags().StringVar(&api, "i", "", "push api docs to yapi")
}
