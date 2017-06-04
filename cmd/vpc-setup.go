package cmd

import "github.com/spf13/cobra"

var vpcSetupCmd = &cobra.Command{
	Use:   "vpc-setup",
	Short: "Provisions VPC, Security Groups and subnets suitable for future development",
	Long: `Provisions VPC, Security Groups and subnets suitable for future development`,

	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.AddCommand(vpcSetupCmd)
}