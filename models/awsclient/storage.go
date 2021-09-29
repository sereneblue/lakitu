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

func (c *AWSClient) CopySnapshot(snapshotId string, sourceRegion string, destRegion string) (string, error) {
	config := c.Config
	config.Region = destRegion

	client := ec2.NewFromConfig(config)

	res, err := client.CopySnapshot(context.TODO(), &ec2.CopySnapshotInput{
		SourceRegion:     aws.String(sourceRegion),
		SourceSnapshotId: aws.String(snapshotId),
	})

	if err == nil {
		// wait for snapshot to be available
		for {
			snapshotState, err := c.GetSnapshotState(*res.SnapshotId, destRegion)

			if err != nil {
				_, deleteErr := c.DeleteSnapshot(*res.SnapshotId, destRegion)

				if deleteErr != nil {
					return "", deleteErr
				}

				return "", err
			}

			if snapshotState == types.SnapshotStateCompleted {
				return *res.SnapshotId, nil
			} else if snapshotState == types.SnapshotStatePending {
				time.Sleep(30 * time.Second)
			} else {
				return "", errors.New("Invalid state for snapshot: " + *res.SnapshotId)
			}
		}
	}

	return "", err
}

func (c *AWSClient) DeleteSnapshot(snapshotId string, region string) (bool, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

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

func (c *AWSClient) GetSnapshotState(snapshotId string, region string) (types.SnapshotState, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DescribeSnapshots(context.TODO(), &ec2.DescribeSnapshotsInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("snapshot-id"),
				Values: []string{snapshotId},
			},
		},
	})

	if err == nil {
		if len(res.Snapshots) > 0 {
			return res.Snapshots[0].State, nil
		}

		return "", errors.New("Could not find snapshot: " + snapshotId)
	}

	return "", err
}
