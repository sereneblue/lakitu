package awsclient

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
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

func (c *AWSClient) GetSpotState(spotRequestId string, region string) (types.SpotInstanceState, *string, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DescribeSpotInstanceRequests(context.TODO(), &ec2.DescribeSpotInstanceRequestsInput{
		SpotInstanceRequestIds: []string{spotRequestId},
	})

	if err == nil {
		if len(res.SpotInstanceRequests) > 0 {
			return res.SpotInstanceRequests[0].State, res.SpotInstanceRequests[0].InstanceId, nil
		}

		return "", nil, errors.New("Could not find spot request: " + spotRequestId)
	}

	return "", nil, err
}

func (c *AWSClient) StartInstance(imageId string, snapshotId string, instanceType types.InstanceType, securityGroupId string, region string, machinePwd string) (string, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	startCmd := fmt.Sprintf(`
		<powershell>
		net user Administrator "%s"
		lakitu-cli mount "%s"
		</powershell>
	`, machinePwd, snapshotId)

	launchSpecs := types.RequestSpotLaunchSpecification{
		ImageId:      aws.String(imageId),
		InstanceType: instanceType,
		SecurityGroupIds: []string{
			securityGroupId,
		},
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(startCmd))),
	}

	if instanceType == types.InstanceTypeG4dnXlarge {
		launchSpecs.BlockDeviceMappings = []types.BlockDeviceMapping{
			types.BlockDeviceMapping{
				DeviceName:  aws.String("xvdca"),
				VirtualName: aws.String("ephemeral0"),
			},
		}
	}

	// create spot instance request
	res, err := client.RequestSpotInstances(context.TODO(), &ec2.RequestSpotInstancesInput{
		AvailabilityZoneGroup: aws.String(region),
		LaunchSpecification:   &launchSpecs,
	})

	if err != nil {
		return "", err
	}

	// check status of spot request
	spotRequest := res.SpotInstanceRequests[0]

	for {
		spotState, instanceId, err := c.GetSpotState(*spotRequest.SpotInstanceRequestId, region)

		if err != nil {
			return "", err
		}

		if spotState == types.SpotInstanceStateActive {
			return *instanceId, nil
		} else if spotState == types.SpotInstanceStateOpen {
			time.Sleep(30 * time.Second)
		} else {
			return "", errors.New("Spot request could not be fulfilled: " + *spotRequest.SpotInstanceRequestId)
		}
	}
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
