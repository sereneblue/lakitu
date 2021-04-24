package models

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSClient struct {
	Config aws.Config
}

type AWSRegion struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var AWS_REGIONS = map[string]string{
	"us-east-2":      "US East (Ohio)",
	"us-east-1":      "US East (N. Virginia)",
	"us-west-1":      "US West (N. California)",
	"us-west-2":      "US West (Oregon)",
	"af-south-1":     "Africa (Cape Town)",
	"ap-east-1":      "Asia Pacific (Hong Kong)",
	"ap-south-1":     "Asia Pacific (Mumbai)",
	"ap-northeast-3": "Asia Pacific (Osaka)",
	"ap-northeast-2": "Asia Pacific (Seoul)",
	"ap-southeast-1": "Asia Pacific (Singapore)",
	"ap-southeast-2": "Asia Pacific (Sydney)",
	"ap-northeast-1": "Asia Pacific (Tokyo)",
	"ca-central-1":   "Canada (Central)",
	"eu-central-1":   "Europe (Frankfurt)",
	"eu-west-1":      "Europe (Ireland)",
	"eu-west-2":      "Europe (London)",
	"eu-south-1":     "Europe (Milan)",
	"eu-west-3":      "Europe (Paris)",
	"eu-north-1":     "Europe (Stockholm)",
	"me-south-1":     "Middle East (Bahrain)",
	"sa-east-1":      "South America (São Paulo)",
}

func NewAWSClient(key, secret, region string) AWSClient {
	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: key, SecretAccessKey: secret, SessionToken: "",
			},
		}))

	c := AWSClient{
		Config: cfg,
	}

	return c
}

func (c *AWSClient) IsValidAWSCredentials() (bool, error) {
	client := sts.NewFromConfig(c.Config)
	_, err := client.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *AWSClient) GetRegions() []AWSRegion {
	var regions []AWSRegion

	client := ec2.NewFromConfig(c.Config)
	output, err := client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})

	if err != nil {
		return regions
	}

	for _, r := range output.Regions {
		regions = append(regions, AWSRegion{
			ID:   *r.RegionName,
			Name: AWS_REGIONS[*r.RegionName],
		})
	}

	return regions
}
