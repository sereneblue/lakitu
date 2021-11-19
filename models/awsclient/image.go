package awsclient

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (c *AWSClient) CreateImage(instanceId string, region string) (string, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.CreateImage(context.TODO(), &ec2.CreateImageInput{
		InstanceId:  aws.String(instanceId),
		Name:        aws.String("lakitu_image"),
		Description: aws.String(""),
		BlockDeviceMappings: []types.BlockDeviceMapping{
			types.BlockDeviceMapping{
				DeviceName: aws.String("xvdh"),
				NoDevice: aws.String(""),
			},
		},
		TagSpecifications: []types.TagSpecification{
			types.TagSpecification{
				ResourceType: types.ResourceTypeImage,
				Tags: []types.Tag{
					types.Tag{
						Key:   aws.String(AWS_TAG_KEY),
						Value: aws.String(""),
					},
				},
			},
		},
	})

	if err == nil {
		// wait for image to be available
		for {
			imageState, err := c.GetImageState(*res.ImageId, region)

			if err != nil {
				deleteErr := c.DeleteImage(*res.ImageId, region)

				if deleteErr != nil {
					return "", deleteErr
				}

				return "", err
			}

			if imageState == types.ImageStateAvailable {
				return *res.ImageId, nil
			} else if imageState == types.ImageStatePending {
				time.Sleep(30 * time.Second)
			} else {
				return "", errors.New("Invalid state for image: " + *res.ImageId)
			}
		}
	}

	return "", err
}

func (c *AWSClient) CopyImage(imageId string, machineUuid string, sourceRegion string, destRegion string) (string, error) {
	config := c.Config
	config.Region = destRegion

	client := ec2.NewFromConfig(config)

	res, err := client.CopyImage(context.TODO(), &ec2.CopyImageInput{
		Name:          aws.String(AWS_TAG_KEY + "_image_for_machine_" + machineUuid),
		SourceRegion:  aws.String(sourceRegion),
		SourceImageId: aws.String(imageId),
	})

	if err == nil {
		// wait for image to be available
		for {
			imageState, err := c.GetImageState(*res.ImageId, destRegion)

			if err != nil {
				deleteErr := c.DeleteImage(*res.ImageId, destRegion)

				if deleteErr != nil {
					return "", deleteErr
				}

				return "", err
			}

			if imageState == types.ImageStateAvailable {
				return *res.ImageId, nil
			} else if imageState == types.ImageStatePending {
				time.Sleep(30 * time.Second)
			} else {
				return "", errors.New("Invalid state for image: " + *res.ImageId)
			}
		}
	}

	return "", err
}

func (c *AWSClient) DeleteImage(imageId string, region string) error {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	_, err := client.DeregisterImage(context.TODO(), &ec2.DeregisterImageInput{
		ImageId: aws.String(imageId),
	})

	// deleting a non-existent ami should return true
	if err != nil {
		if !(strings.Contains(err.Error(), "InvalidAMIID.NotFound") ||
			strings.Contains(err.Error(), "InvalidAMIID.Unavailable")) {
			return err
		}
	}

	// delete snapshots associated with ami
	res, err := client.DescribeSnapshots(context.TODO(), &ec2.DescribeSnapshotsInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("description"),
				Values: []string{"Created by*" + imageId + "*"},
			},
		},
	})

	if err != nil {
		return err
	}

	for _, snapshot := range res.Snapshots {
		err := c.DeleteSnapshot(*snapshot.SnapshotId, region)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *AWSClient) GetImageState(imageId string, region string) (types.ImageState, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DescribeImages(context.TODO(), &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("image-id"),
				Values: []string{imageId},
			},
		},
	})

	if err == nil {
		if len(res.Images) > 0 {
			return res.Images[0].State, nil
		}

		return "", errors.New("Could not find image: " + imageId)
	}

	return "", err
}
