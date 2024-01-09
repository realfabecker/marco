package toolkit

import (
	"fmt"
	"github.com/rafaelbeecker/mwskit/internal/mws"
	"github.com/spf13/cobra"
	"log"
)

func newReportCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create report",
		RunE: func(cmd *cobra.Command, args []string) error {
			reportType, _ := cmd.Flags().GetString("report-type")
			s := mws.ReportService{}
			reportId, err := s.CreateReport(reportType)
			if err != nil {
				return err
			}
			log.Printf("report created: %v", reportId)
			return nil
		},
	}
	cmd.Flags().String("report-type", "", "report type to be created")
	_ = cmd.MarkFlagRequired("report-type")
	return cmd
}

func newReportDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download report",
		RunE: func(cmd *cobra.Command, args []string) error {
			target, _ := cmd.Flags().GetString("target")
			reportId, _ := cmd.Flags().GetString("report-id")
			s := mws.ReportService{}
			documentId, err := s.GetReport(reportId)
			if err != nil {
				return err
			}
			if documentId == "" {
				return fmt.Errorf("documentId not ready to consume yet")
			}
			documentUrl, err := s.GetReportDocument(documentId)
			if err != nil {
				return err
			}
			if err := s.DownloadReportDocument(reportId, target, documentUrl); err != nil {
				return err
			}
			fmt.Printf("Download de arquivo relatório %s no diretório %s\n", reportId, target)
			return nil
		},
	}
	cmd.Flags().String("target", "", "report destination directory")
	cmd.Flags().String("report-id", "", "report id to download")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("report-id")
	return cmd
}

func newReportCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "report",
		Short: "Report toolset",
	}
	root.AddCommand(newReportDownloadCmd())
	root.AddCommand(newReportCreateCmd())
	return root
}
