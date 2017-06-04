package cmd

import "github.com/spf13/cobra"

var iamScanCmd = &cobra.Command{
	Use:   "iam-scan",
	Short: "Lists all IAM users with associated rules, detects too wide permissions and unused accounts/roles/policies",
	Long: `Lists all IAM users with associated rules, detects too wide permissions and unused accounts/roles/policies`,

	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.AddCommand(iamScanCmd)
}