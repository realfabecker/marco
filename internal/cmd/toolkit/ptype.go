package toolkit

import (
	"log"
	"path/filepath"

	"github.com/rafaelbeecker/mwskit/internal/mws"
	"github.com/spf13/cobra"
)

func newPtypeSchemaDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download Product Type schema definitions",
		RunE: func(cmd *cobra.Command, args []string) error {
			productType, _ := cmd.Flags().GetString("ptype")
			productList, _ := cmd.Flags().GetString("plist")
			target, _ := cmd.Flags().GetString("target")
			seller, _ := cmd.Flags().GetString("seller")

			if productType != "" {
				s := mws.PtypeService{}
				link, err := s.GetProductTypeDefSchemaUrl(seller, productType)
				if err != nil {
					return err
				}
				dest := filepath.Join(target, productType+".json")
				if err := s.DownloadProductTypeDef(dest, link); err != nil {
					return err
				}
				log.Printf("schema downloaded at %s\n", dest)
			} else if productList != "" {
				s := mws.PtypeService{}
				if err := s.DownloadBatchTypeDef(
					seller,
					productList,
					target,
				); err != nil {
					return err
				}
				log.Println("batch downloaded successfuly")
			}
			return nil
		},
	}

	cmd.Flags().String("seller", "", "seller")
	cmd.Flags().String("ptype", "", "product type")
	cmd.Flags().String("target", "", "download target")
	cmd.Flags().String("plist", "", "product type list (csv)")

	cmd.MarkFlagRequired("seller")
	cmd.MarkFlagRequired("target")
	return cmd
}

func newPtypeCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "ptype",
		Short: "Product Type toolset",
	}
	root.AddCommand(newPtypeSchemaDownloadCmd())
	return root
}
