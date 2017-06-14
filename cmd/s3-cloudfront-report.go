package cmd

import "github.com/spf13/cobra"

var s3CloudfrontReportCmd = &cobra.Command{
	Use:   "s3-cloudfront-report",
	Short: "List S3 buckets with connected Cloudfront distributions sorted by size.",
	Long: `List S3 buckets with connected Cloudfront distributions sorted by size.`,

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(s3CloudfrontReportCmd)
}