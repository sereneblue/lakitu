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

func (c *AWSClient) checkSnapshotState(snapshotId string, region string) (string, error) {
	for {
		snapshotState, err := c.GetSnapshotState(snapshotId, region)

		if err != nil {
			deleteErr := c.DeleteSnapshot(snapshotId, region)

			if deleteErr != nil {
				return "", deleteErr
			}

			return "", err
		}

		if snapshotState == types.SnapshotStateCompleted {
			return snapshotId, nil
		} else if snapshotState == types.SnapshotStatePending {
			time.Sleep(30 * time.Second)
		} else {
			return "", errors.New("Invalid state for snapshot: " + snapshotId)
		}
	}
}

func (c *AWSClient) CreateSnapshot(volumeId string, region string) (string, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.CreateSnapshot(context.TODO(), &ec2.CreateSnapshotInput{
		VolumeId: &volumeId,
	})

	if err == nil {
		return c.checkSnapshotState(*res.SnapshotId, region)
	}

	return "", err
}

func (c *AWSClient) CreateNewVolume(instanceId string, size int32, region string) error {
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

	if err != nil {
		return err
	}

	if len(res.Reservations[0].Instances) == 0 {
		return errors.New("Could not find instance info")
	}

	volumeRes, err := client.CreateVolume(context.TODO(), &ec2.CreateVolumeInput{
		AvailabilityZone: res.Reservations[0].Instances[0].Placement.AvailabilityZone,
		Size:             size,
		VolumeType:       types.VolumeTypeGp3,
	})

	if err == nil {
		// wait for volume to be available
		for {
			volumeState, err := c.GetVolumeState(*volumeRes.VolumeId, region)

			if err != nil {
				_, deleteErr := c.DeleteVolume(*volumeRes.VolumeId, region)

				if deleteErr != nil {
					return deleteErr
				}

				return err
			}

			if volumeState == types.VolumeStateAvailable {
				break
			} else if volumeState == types.VolumeStateCreating {
				time.Sleep(30 * time.Second)
			} else {
				return errors.New("Invalid state for volume: " + *volumeRes.VolumeId)
			}
		}
	} else {
		return err
	}

	// attach volume to instance
	_, err = client.AttachVolume(context.TODO(), &ec2.AttachVolumeInput{
		Device: aws.String("xvdh"),
		InstanceId: aws.String(instanceId),
		VolumeId:   volumeRes.VolumeId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *AWSClient) CreateVolume(snapshotId string, region string) (string, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DescribeAvailabilityZones(context.TODO(), &ec2.DescribeAvailabilityZonesInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("region-name"),
				Values: []string{region},
			},
			types.Filter{
				Name:   aws.String("state"),
				Values: []string{"available"},
			},
		},
	})

	if err != nil {
		return "", err
	}

	if len(res.AvailabilityZones) == 0 {
		return "", errors.New("Could not find availability zones")
	}

	volumeRes, err := client.CreateVolume(context.TODO(), &ec2.CreateVolumeInput{
		AvailabilityZone: res.AvailabilityZones[0].ZoneName,
		SnapshotId:       aws.String(snapshotId),
		VolumeType:       types.VolumeTypeGp3,
	})

	if err == nil {
		// wait for volume to be available
		for {
			volumeState, err := c.GetVolumeState(*volumeRes.VolumeId, region)

			if err != nil {
				_, deleteErr := c.DeleteVolume(*volumeRes.VolumeId, region)

				if deleteErr != nil {
					return "", deleteErr
				}

				return "", err
			}

			if volumeState == types.VolumeStateAvailable {
				return *volumeRes.VolumeId, nil
			} else if volumeState == types.VolumeStateCreating {
				time.Sleep(30 * time.Second)
			} else {
				return "", errors.New("Invalid state for volume: " + *volumeRes.VolumeId)
			}
		}
	}

	return "", err
}

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
		return c.checkSnapshotState(*res.SnapshotId, destRegion)
	}

	return "", err
}

func (c *AWSClient) DeleteSnapshot(snapshotId string, region string) error {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DeleteSnapshot(context.TODO(), &ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshotId),
	})

	if res != nil {
		return nil
	}

	// deleting a non-existent snapshot should return true
	if err != nil {
		if strings.Contains(err.Error(), "InvalidSnapshot.NotFound") {
			return nil
		}
	}

	return err
}

func (c *AWSClient) DeleteVolume(volumeId string, region string) (bool, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DeleteVolume(context.TODO(), &ec2.DeleteVolumeInput{
		VolumeId: aws.String(volumeId),
	})

	if res != nil {
		return true, nil
	}

	// deleting a non-existent volume should return true
	if err != nil {
		if strings.Contains(err.Error(), "InvalidVolume.NotFound") {
			return true, nil
		}
	}

	return false, err
}

func (c *AWSClient) GetVolumeState(volumeId string, region string) (types.VolumeState, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DescribeVolumes(context.TODO(), &ec2.DescribeVolumesInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("volume-id"),
				Values: []string{volumeId},
			},
		},
	})

	if err == nil {
		if len(res.Volumes) > 0 {
			return res.Volumes[0].State, nil
		}

		return "", errors.New("Could not find volume: " + volumeId)
	}

	return "", err
}

func (c *AWSClient) GetVolumeModificationState(volumeId string, region string) (types.VolumeModificationState, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	res, err := client.DescribeVolumesModifications(context.TODO(), &ec2.DescribeVolumesModificationsInput{
		Filters: []types.Filter{
			types.Filter{
				Name:   aws.String("volume-id"),
				Values: []string{volumeId},
			},
		},
	})

	if err == nil {
		if len(res.VolumesModifications) > 0 {
			return res.VolumesModifications[0].ModificationState, nil
		}

		return "", errors.New("Could not find volume: " + volumeId)
	}

	return "", err
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

func (c *AWSClient) ModifyVolume(volumeId string, newSize int32, region string) (bool, error) {
	config := c.Config
	config.Region = region

	client := ec2.NewFromConfig(config)

	_, err := client.ModifyVolume(context.TODO(), &ec2.ModifyVolumeInput{
		VolumeId:   aws.String(volumeId),
		Size:       newSize,
		VolumeType: types.VolumeTypeGp3,
	})

	if err == nil {
		// wait for volume modification status to be available
		for {
			volumeModState, err := c.GetVolumeModificationState(volumeId, region)

			if err != nil {
				_, deleteErr := c.DeleteVolume(volumeId, region)

				if deleteErr != nil {
					return false, deleteErr
				}

				return false, err
			}

			if volumeModState == types.VolumeModificationStateCompleted {
				return true, nil
			} else if volumeModState == types.VolumeModificationStateModifying ||
				volumeModState == types.VolumeModificationStateOptimizing {
				time.Sleep(30 * time.Second)
			} else {
				return false, errors.New("Invalid state for modification to volume: " + volumeId)
			}
		}
	} else {
		return false, err
	}
}

// Resizes snapshot, doesn't automate expanding partition
// Expanding partition will occur on next machine boot
func (c *AWSClient) ResizeSnapshot(snapshotId string, newSize int32, region string) (string, error) {
	volumeId, err := c.CreateVolume(snapshotId, region)
	if err != nil {
		return "", err
	}

	_, err = c.ModifyVolume(volumeId, newSize, region)
	if err != nil {
		return "", err
	}

	newSnapshotId, err := c.CreateSnapshot(volumeId, region)
	if err != nil {
		return "", err
	}

	_, err = c.DeleteVolume(volumeId, region)
	if err != nil {
		return "", err
	}

	err = c.DeleteSnapshot(snapshotId, region)
	if err != nil {
		return "", err
	}

	return newSnapshotId, nil
}
