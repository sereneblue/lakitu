package awsclient

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type AWSSecurityGroup struct {
	Id             int64
	GroupId        string
	StreamSoftware string
	Created        time.Time `xorm:"created"`
}

func (sg *AWSSecurityGroup) TableName() string {
	return "aws_securitygroups"
}

func (c *AWSClient) CreateSecurityGroup(streamSW StreamSoftware) (AWSSecurityGroup, error) {
	var group AWSSecurityGroup

	client := ec2.NewFromConfig(c.Config)

	sg, err := client.CreateSecurityGroup(context.TODO(), &ec2.CreateSecurityGroupInput{
		Description: aws.String("Cloud gaming security group"),
		GroupName:   aws.String("lakitu security group - " + time.Now().Format("20060102150405")),
		TagSpecifications: []types.TagSpecification{
			types.TagSpecification{
				ResourceType: types.ResourceTypeSecurityGroup,
				Tags: []types.Tag{
					types.Tag{
						Key:   aws.String(AWS_TAG_KEY),
						Value: aws.String(""),
					},
				},
			},
		},
	})

	if err != nil {
		return group, err
	}

	if streamSW == PARSEC {
		_, err = client.AuthorizeSecurityGroupIngress(context.TODO(), &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: sg.GroupId,
			IpPermissions: []types.IpPermission{
				types.IpPermission{
					FromPort:   8000,
					ToPort:     8200,
					IpProtocol: aws.String("udp"),
					IpRanges: []types.IpRange{
						types.IpRange{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				types.IpPermission{
					FromPort:   9000,
					ToPort:     9200,
					IpProtocol: aws.String("udp"),
					IpRanges: []types.IpRange{
						types.IpRange{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
			},
		})

		if err != nil {
			return group, err
		}
	} else {
		// moonlight
		_, err = client.AuthorizeSecurityGroupIngress(context.TODO(), &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: sg.GroupId,
			IpPermissions: []types.IpPermission{
				types.IpPermission{
					FromPort:   47984,
					ToPort:     47984,
					IpProtocol: aws.String("tcp"),
					IpRanges: []types.IpRange{
						types.IpRange{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				types.IpPermission{
					FromPort:   47989,
					ToPort:     47989,
					IpProtocol: aws.String("tcp"),
					IpRanges: []types.IpRange{
						types.IpRange{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				types.IpPermission{
					FromPort:   48010,
					ToPort:     48010,
					IpProtocol: aws.String("tcp"),
					IpRanges: []types.IpRange{
						types.IpRange{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				types.IpPermission{
					FromPort:   47998,
					ToPort:     48000,
					IpProtocol: aws.String("udp"),
					IpRanges: []types.IpRange{
						types.IpRange{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				types.IpPermission{
					FromPort:   48010,
					ToPort:     48010,
					IpProtocol: aws.String("udp"),
					IpRanges: []types.IpRange{
						types.IpRange{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
				types.IpPermission{
					FromPort:   3389,
					ToPort:     3389,
					IpProtocol: aws.String("tcp"),
					IpRanges: []types.IpRange{
						types.IpRange{
							CidrIp: aws.String("0.0.0.0/0"),
						},
					},
				},
			},
		})

		if err != nil {
			return group, err
		}
	}

	group.GroupId = *sg.GroupId
	group.StreamSoftware = streamSW.String()

	return group, nil
}

func (c *AWSClient) GetSecurityGroups() ([]types.SecurityGroup, error) {
	client := ec2.NewFromConfig(c.Config)
	output, err := client.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("tag-key"),
				Values: []string{AWS_TAG_KEY},
			},
		},
	})

	if err != nil {
		return []types.SecurityGroup{}, err
	}

	return output.SecurityGroups, nil
}
