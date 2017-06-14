package cmd

import "github.com/spf13/cobra"

var disasterRecoveryCmd = &cobra.Command{
	Use:   "disaster-recovery",
	Short: "Provisions EC2 instance and RDS instance basing on last snapshot provided.",
	Long: `Provisions EC2 instance and RDS instance basing on last snapshot provided.`,

	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.AddCommand(disasterRecoveryCmd)
}