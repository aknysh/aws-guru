package cmd

/* Shared variables */
var region string
var profile string
var cfgFile string

/* ec2-snapshoter variables */
var cronPattern string
var cronName string
var reattachOnly bool

/* iam-aim-linker variables */
var accountId string
var path string
var slaveProfile string
var slaveName string
var format string

/* vpc-setup variables */
var vpcCidr string
var privateSubnetCidr string
var publicSubnetCidr string
