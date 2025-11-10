package internationalization

import (
	"github.com/spf13/cobra"
)

var (
	// Cmd describes a kube command.
	Cmd = &cobra.Command{
		Use:   "lang",
		Short: "Generate international files",
	}

	multilingualCmd = &cobra.Command{
		Use:   "multiInit",
		Short: "Generate international language package",
		RunE:  MultilingualAction,
	}
)

func init() {
	multilingualCmd.Flags().StringVar(&OwnerString, "owner", "", "The github owner")
	multilingualCmd.Flags().StringVar(&RepoString, "repo", "", "The github owner repo")
	multilingualCmd.Flags().StringVar(&PathString, "path", "", "The github owner repo path")
	multilingualCmd.Flags().StringVar(&MovePathString, "movePath", "", "Move github owner repo path to local path")

	Cmd.AddCommand(multilingualCmd)
}
