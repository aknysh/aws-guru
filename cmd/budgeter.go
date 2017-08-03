package cmd

import "github.com/spf13/cobra"

var budgeterCmd = &cobra.Command{
	Use:   "budgeter",
	Short: "Creates budgets for EC2, RDS, S3 & Cloudfront usage.",
	Long:  `Creates budgets for EC2, RDS, S3 & Cloudfront usage.`,

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(budgeterCmd)
}
