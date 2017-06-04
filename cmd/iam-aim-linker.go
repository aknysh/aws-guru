package cmd

import "github.com/spf13/cobra"

var iamAimLinkerCmd = &cobra.Command{
	Use:   "iam-scan",
	Short: "Links account to AWS-IAM-Manager",
	Long: `Links account to AWS-IAM-Manager`,

	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.AddCommand(iamAimLinkerCmd)
}