package toolkit

import (
	"github.com/spf13/cobra"
)

// NewRootCmd
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:  "tookit",
		Long: "AWS Selling Partner API toolkit",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	root.SetHelpCommand(&cobra.Command{Hidden: true})
	root.AddCommand(newBtreeCmd())
	root.AddCommand(newPtypeCmd())
	root.AddCommand(newReportCmd())
	return root
}
