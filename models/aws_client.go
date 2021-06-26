package models

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	pricingTypes "github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSClient struct {
	Config aws.Config
}

type AWSGPUInstance struct {
	InstanceType string  `json:"instance"`
	Price        float64 `json:"price"`
}

type AWSPrices struct {
	Bandwidth float64 `json:"bandwidth"`
	Volume    float64 `json:"volume"`
	Snapshots float64 `json:"snapshot"`
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
	"eu-central-1":   "EU (Frankfurt)",
	"eu-west-1":      "EU (Ireland)",
	"eu-west-2":      "EU (London)",
	"eu-south-1":     "EU (Milan)",
	"eu-west-3":      "EU (Paris)",
	"eu-north-1":     "EU (Stockholm)",
	"me-south-1":     "Middle East (Bahrain)",
	"sa-east-1":      "South America (Sao Paulo)",
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

func (c *AWSClient) GetGPUInstances(region string) []AWSGPUInstance {
	var instances []AWSGPUInstance

	c.Config.Region = region

	instancePriceMap := map[string]float64{}

	client := ec2.NewFromConfig(c.Config)
	timestamp := time.Now()
	output, err := client.DescribeSpotPriceHistory(context.TODO(), &ec2.DescribeSpotPriceHistoryInput{
		InstanceTypes: []types.InstanceType{
			types.InstanceTypeG3sXlarge,
			types.InstanceTypeG4dnXlarge,
		},
		ProductDescriptions: []string{"Windows"},
		StartTime:           &timestamp,
	})

	if err != nil {
		return instances
	}

	for _, sp := range output.SpotPriceHistory {
		instanceType := string(sp.InstanceType)
		instancePrice, _ := strconv.ParseFloat(*sp.SpotPrice, 64)

		if price, ok := instancePriceMap[instanceType]; ok {
			avgPrice := (price + instancePrice) / 2
			instancePriceMap[instanceType] = avgPrice
		} else {
			instancePriceMap[instanceType] = instancePrice
		}
	}

	for k, v := range instancePriceMap {
		instances = append(instances, AWSGPUInstance{
			InstanceType: k,
			Price:        v,
		})
	}

	return instances
}

func (c *AWSClient) GetPrices(region string) AWSPrices {
	prices := AWSPrices{}
	re := regexp.MustCompile(`"pricePerUnit":{"USD":"(.*?)"}`)

	client := pricing.NewFromConfig(c.Config)
	resGP3, err := client.GetProducts(context.TODO(), &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []pricingTypes.Filter{
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("productFamily"),
				Value: aws.String("Storage"),
			},
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("volumeApiName"),
				Value: aws.String("gp3"),
			},
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("location"),
				Value: aws.String(AWS_REGIONS[region]),
			},
		},
	})

	if err == nil {
		if len(resGP3.PriceList) > 0 {
			storagePrice := re.FindStringSubmatch(resGP3.PriceList[0])
			prices.Volume, _ = strconv.ParseFloat(storagePrice[1], 64)
		} else {
			resGP2, err := client.GetProducts(context.TODO(), &pricing.GetProductsInput{
				ServiceCode: aws.String("AmazonEC2"),
				Filters: []pricingTypes.Filter{
					pricingTypes.Filter{
						Type:  pricingTypes.FilterTypeTermMatch,
						Field: aws.String("productFamily"),
						Value: aws.String("Storage"),
					},
					pricingTypes.Filter{
						Type:  pricingTypes.FilterTypeTermMatch,
						Field: aws.String("volumeApiName"),
						Value: aws.String("gp2"),
					},
					pricingTypes.Filter{
						Type:  pricingTypes.FilterTypeTermMatch,
						Field: aws.String("location"),
						Value: aws.String(AWS_REGIONS[region]),
					},
				},
			})

			if err == nil {
				if len(resGP2.PriceList) > 0 {
					storagePrice := re.FindStringSubmatch(resGP2.PriceList[0])
					prices.Volume, _ = strconv.ParseFloat(storagePrice[1], 64)
				}
			}
		}
	}

	resSnapShot, err := client.GetProducts(context.TODO(), &pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters: []pricingTypes.Filter{
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("productFamily"),
				Value: aws.String("Storage Snapshot"),
			},
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("storageMedia"),
				Value: aws.String("Amazon S3"),
			},
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("location"),
				Value: aws.String(AWS_REGIONS[region]),
			},
		},
	})

	if err == nil {
		snapShotPrice := re.FindStringSubmatch(resSnapShot.PriceList[0])
		prices.Snapshots, _ = strconv.ParseFloat(snapShotPrice[1], 64)
	}

	resBandwidth, err := client.GetProducts(context.TODO(), &pricing.GetProductsInput{
		ServiceCode: aws.String("AWSDataTransfer"),
		Filters: []pricingTypes.Filter{
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("toLocationType"),
				Value: aws.String("Other"),
			},
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("toLocation"),
				Value: aws.String("External"),
			},
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("transferType"),
				Value: aws.String("AWS Outbound"),
			},
			pricingTypes.Filter{
				Type:  pricingTypes.FilterTypeTermMatch,
				Field: aws.String("fromLocation"),
				Value: aws.String(AWS_REGIONS[region]),
			},
		},
	})

	if err == nil {
		reBW := regexp.MustCompile(`("endRange":"10240".*?"beginRange":"1".*?}})`)

		bandwidthPrice := reBW.FindStringSubmatch(resBandwidth.PriceList[0])
		bandwidthPrice = re.FindStringSubmatch(bandwidthPrice[1])
		prices.Bandwidth, _ = strconv.ParseFloat(bandwidthPrice[1], 64)
	}

	return prices
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
