# aws-guru

### Goal
Set of provisioning and maintenance scripts for easier AWS accounts management.


### Design decisions
Tool is going to have two types of scripts:

- Initial Setup & provisioning scripts - goal of these scripts is to prepare clean/semi-clean AWS account for work - enable Cloudtrail, create budgets and billing dashboards, check IAM roles, setup automatic volumes backups etc., more on that: https://www.netguru.co/blog/my-first-5-minutes-on-an-aws-account
- Maintenance scripts - mostly check scripts providing information about possible configuration drift, for instance: Are S3 assets served through Cloudfront, Are RDS backups present, are there any cloudwatch alarms in place,

Tool is going to be written in Go because:
- Go AWS SDK well documented
- It’s fast
- imperative control flow is very easy to handle and fits perfectly for such needs as calling multiple APIs sequentially
- Provides nice API for concurrency
- Can be used as AWS Lambda using Node.js shim

Every script should follow single responsibility principle and should be a separate file so tool could be easily extendable.

### List of Scripts

 Script | Description | Maintenance or Provision | Implementation status
 --- | --- | --- | ---
 s3-cloudfront-report | List S3 buckets with connected Cloudfront distributions sorted by size. | M | Done
 ec2-snapshoter | Sets up CloudWatch scheduled expression to take snaphot of all ec2 volumes every 24h | M/P | Done
 iam-scan | Lists all IAM users with associated rules, detects too wide permissions and unused accounts/roles/policies | M | To do
 iam-repair | Fixes issues found by `iam-scan` | M | To do
 iam-aim-linker | Links account to AWS-IAM-Manager | P | To do
 iam-setup | Provides “backdoor” IAM user with admin access for devops, enforces “password policy”, deletes root access keys, etc | P | To do
 ec2-scan | Lists all EC2 instances with associated security group, detects anomalies e.g. too wide access to instances or lack of termination protection | M | To do
 vpc-setup | Provisions VPC, Security Groups and subnets suitable for future development | P | In progress
 budgeter | Creates budgets for EC2, RDS, S3 & Cloudfront usage. | P | To do
 disaster-recovery | Provisions EC2 instance and RDS instance basing on last snapshot provided | M | To do

### Development

- Clone this repository to `$GOPATH/src/netguru/aws-guru`
- Run `dep ensure` to download all the dependencies. If you don't have `dep` you can get it here: https://github.com/golang/dep
- Build & Run the app `go run main.go [arg]`
