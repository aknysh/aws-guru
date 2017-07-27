package cmd

import (
	"fmt"
	awslib "netguru/aws-guru/aws"
	"netguru/aws-guru/utils"

	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/spf13/cobra"

	"code.cloudfoundry.org/bytefmt"
)

var s3CloudfrontReportCmd = &cobra.Command{
	Use:   "s3-cloudfront-report",
	Short: "List S3 buckets with connected Cloudfront distributions sorted by size.",
	Long:  `List S3 buckets with connected Cloudfront distributions sorted by size.`,

	Run: func(cmd *cobra.Command, args []string) {
		s3CloudfrontReportRun()
	},
}

type bucketDistro struct {
	Bucket   string
	CfDistro string
}

type s3CloudfrontReportJSON struct {
	Bucket    string `json:"bucket"`
	CfURL     string `json:"cf_url"`
	Region    string `json:"region"`
	Size      int64  `json:"size"`
	SizeHuman string `json:"size_human"`
}

func init() {
	s3CloudfrontReportCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "region")
	s3CloudfrontReportCmd.Flags().StringVarP(&profile, "profile", "p", "default", "profile to use with credentials")
	s3CloudfrontReportCmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table/json)")
	RootCmd.AddCommand(s3CloudfrontReportCmd)
}

func s3CloudfrontReportRun() {
	sess := awslib.CreateProfiledSession(&region, &profile)
	cfSvc := awslib.CreateCloudfrontContext(sess)

	fmt.Printf("Looking for CF distributions connected with s3... ")

	distro, err := awslib.ListDistributions(cfSvc)
	if err != nil {
		utils.ExitWithError(err)
	}

	var buckets []bucketDistro

	for _, item := range distro.Items {
		for _, origin := range item.Origins.Items {
			if strings.Contains(*origin.DomainName, ".s3.amazonaws.com") {
				buckets = append(buckets, bucketDistro{Bucket: *origin.DomainName, CfDistro: *item.DomainName})
			}
		}
	}

	if len(buckets) > 0 {
		fmt.Printf("Found %d buckets\n", len(buckets))
		calculateBucketSizes(buckets)
	} else {
		fmt.Printf("No buckets found\n")
	}
}

func fetchBucketName(domainName string) string {
	splitted := strings.Split(domainName, ".s3.amazonaws.com")
	return splitted[0]
}

func fetchBucketRegion(bucket string, sess *session.Session) string {
	ctx := context.Background()

	bucketRegion, err := s3manager.GetBucketRegion(ctx, sess, bucket, region)
	if err != nil {
		utils.ExitWithError(err)
	}

	return bucketRegion
}

func calculateBucketSizes(buckets []bucketDistro) {
	sess := awslib.CreateProfiledSession(&region, &profile)
	output := []s3CloudfrontReportJSON{}

	for _, bd := range buckets {
		distro := bd.CfDistro
		bucket := fetchBucketName(bd.Bucket)

		bucketRegion := fetchBucketRegion(bucket, sess)

		sess = awslib.CreateProfiledSession(&bucketRegion, &profile)
		s3Svc := awslib.CreateS3Context(sess)

		input := &s3.ListObjectsV2Input{
			Bucket: &bucket,
		}

		var bucketSize int64

		err := s3Svc.ListObjectsV2Pages(input, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			fmt.Print(".")
			for _, obj := range page.Contents {
				bucketSize = bucketSize + *obj.Size
			}
			return lastPage
		})

		if err != nil {
			utils.ExitWithError(err)
		}

		output = append(output, s3CloudfrontReportJSON{
			Bucket:    bucket,
			CfURL:     distro,
			Region:    bucketRegion,
			Size:      bucketSize,
			SizeHuman: bytefmt.ByteSize(uint64(bucketSize)),
		})
	}

	fmt.Print("\n\n")

	switch format {
	case "json":
		utils.PrintJSON(output)
	default:
		utils.PrintTable(s3CloudfrontReportToTable(output), []string{"Bucket", "CF URL", "Region", "Size"})
	}

}

func s3CloudfrontReportToTable(json []s3CloudfrontReportJSON) [][]string {
	table := [][]string{}

	for _, item := range json {
		table = append(table, []string{item.Bucket, item.CfURL, item.Region, item.SizeHuman})
	}

	return table
}
