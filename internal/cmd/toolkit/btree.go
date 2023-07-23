package toolkit

import (
	"fmt"
	"os"

	"github.com/rafaelbeecker/mwskit/internal/mws"
	"github.com/spf13/cobra"
)

func newBtreeReportFlatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flat",
		Short: "Flat xml report browse tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, _ := cmd.Flags().GetString("report")
			target, _ := cmd.Flags().GetString("target")

			if _, err := os.Stat(target); err != nil {
				return fmt.Errorf("open-target: %w", err)
			}

			if _, err := os.Stat(report); err != nil {
				return fmt.Errorf("open-report: %w", err)
			}

			s := mws.BrowseService{}
			l, err := s.Read(report)
			if err != nil {
				return err
			}
			return s.Flat(l, target)
		},
	}

	cmd.Flags().String("report", "", "product type")
	cmd.Flags().String("target", "", "download target")

	cmd.MarkFlagRequired("report")
	cmd.MarkFlagRequired("target")
	return cmd
}

func newBtreeCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "btree",
		Short: "Browse Tree toolset",
	}
	root.AddCommand(newBtreeReportFlatCmd())
	return root
}
