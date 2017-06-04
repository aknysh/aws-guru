package cmd

import "github.com/spf13/cobra"

var iamSetupCmd = &cobra.Command{
	Use:   "iam-setup",
	Short: "Provides “backdoor” IAM user with admin access for devops, enforces “password policy”, deletes root access keys, etc",
	Long: `Provides “backdoor” IAM user with admin access for devops, enforces “password policy”, deletes root access keys, etc`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.AddCommand(iamSetupCmd)
}