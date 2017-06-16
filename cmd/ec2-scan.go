package cmd

import (
	"github.com/spf13/cobra"
	"github.com/olekukonko/tablewriter"
)

var ec2ScanCmd = &cobra.Command{
	Use:   "ec2-scan",
	Short: "Lists all EC2 instances with associated security group, detects anomalies e.g. too wide access to instances or lack of termination protection",
	Long: `Lists all EC2 instances with associated security group, detects anomalies e.g. too wide access to instances or lack of termination protection`,

	Run: func(cmd *cobra.Command, args []string) {
		ec2ScanRun()
	},
}

func init() {
	RootCmd.AddCommand(ec2ScanCmd)
}

func ec2ScanRun() {

}