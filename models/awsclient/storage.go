package awsclient

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (c *AWSClient) DeleteImage(imageId string) (bool, error) {
	client := ec2.NewFromConfig(c.Config)

	_, err := client.DeregisterImage(context.TODO(), &ec2.DeregisterImageInput{
		ImageId: aws.String(imageId),
	})

	// deleting a non-existent ami should return true
	if err != nil {
		if !(strings.Contains(err.Error(), "InvalidAMIID.NotFound") || 
		   strings.Contains(err.Error(), "InvalidAMIID.Unavailable")) {
			return false, err
		}
	}

	// delete snapshots associated with ami
	res, err := client.DescribeSnapshots(context.TODO(), &ec2.DescribeSnapshotsInput{
		Filters: []types.Filter{
			types.Filter{
				Name: aws.String("description"),
				Values: []string{"Created by*" + imageId + "*"},
			},
		},
	})

	if err != nil {
		return false, err
	}

	for _, snapshot := range res.Snapshots {
		_, err := c.DeleteSnapshot(*snapshot.SnapshotId)

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (c *AWSClient) DeleteSnapshot(snapshotId string) (bool, error) {
	client := ec2.NewFromConfig(c.Config)

	res, err := client.DeleteSnapshot(context.TODO(), &ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshotId),
	})

	if res != nil {
		return true, nil
	}

	// deleting a non-existent snapshot should return true
	if err != nil {
		if strings.Contains(err.Error(), "InvalidSnapshot.NotFound") {
			return true, nil
		}
	}

	return false, err
}