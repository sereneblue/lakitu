package awsclient

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (c *AWSClient) CreateSpotInstance(imageId string, createMachine bool) (bool, error) {
	client := ec2.NewFromConfig(c.Config)

	_, err := client.RequestSpotInstances(context.TODO(), &ec2.RequestSpotInstancesInput{
		ClientToken:   aws.String(""),
		InstanceCount: 1,
		LaunchSpecification: &types.RequestSpotLaunchSpecification{
			ImageId:          aws.String(""),
			InstanceType:     types.InstanceTypeG22xlarge,
			KeyName:          aws.String(""),
			SecurityGroupIds: []string{},
			UserData:         aws.String(""),
		},
		SpotPrice: aws.String(""),
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
		return false, err
	}

	return true, nil
}

func (c *AWSClient) GetInstanceState(instanceId string, region string) (types.InstanceStateName, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("instance-id"),
				Values: []string{instanceId},
			},
		},
	})

	if err == nil {
		if len(res.Reservations) > 0 {
			return res.Reservations[0].Instances[0].State.Name, nil
		}

		return "", errors.New("Could not find instance: " + instanceId)
	}

	return "", err
}

func (c *AWSClient) TerminateInstance(instanceId string, region string) error {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceId},
	})

	if err != nil {
		return err
	}

	if len(res.TerminatingInstances) == 0 {
		return errors.New("Missing instance in termination list")
	}

	// wait for instance to be terminated
	for {
		instanceState, err := c.GetInstanceState(instanceId, region)

		if err != nil {
			return err
		}

		if instanceState == types.InstanceStateNameTerminated {
			return nil
		} else if instanceState == types.InstanceStateNameShuttingDown {
			time.Sleep(30 * time.Second)
		} else {
			return errors.New("Invalid state for instance: " + instanceId)
		}
	}
}
