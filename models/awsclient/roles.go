package awsclient

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type AWSRole struct {
	Id       int64
	RoleId   string
	RoleName string
	Arn      string
	Created  time.Time `xorm:"created"`
}

func (a *AWSRole) TableName() string {
	return "aws_roles"
}

func (c *AWSClient) CreateRole() (AWSRole, error) {
	var role AWSRole

	client := iam.NewFromConfig(c.Config)

	roleName := "lakitu_role_" + time.Now().Format("20060102150405")
	output, err := client.CreateRole(context.TODO(), &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(`{
			  "Version": "2012-10-17",
			  "Statement": [
			    {
			      "Effect": "Allow",
			      "Principal": {
			        "Service": "ec2.amazonaws.com"
			      },
			      "Action": "sts:AssumeRole"
			    }
			  ]
			}
		`),
		MaxSessionDuration: aws.Int32(1 * 60 * 60),
		RoleName:           aws.String(roleName),
		Tags: []types.Tag{
			types.Tag{
				Key:   aws.String(AWS_TAG_KEY),
				Value: aws.String(""),
			},
		},
	})

	if err != nil {
		return role, err
	}

	_, err = client.PutRolePolicy(context.TODO(), &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(`{
		    "Version": "2012-10-17",
		    "Statement": [
		        {
					"Sid": "LakituManageStorage",
		            "Action": [
			            "ec2:AttachVolume",
		            	"ec2:CreateVolume",
			            "ec2:DescribeVolumes"
		            ],
		            "Effect": "Allow",
		            "Resource": [
		            	"*"
		            ]
		        },
		    ]
		}`),
		PolicyName: aws.String("EC2_Limited_Access"),
		RoleName:   aws.String(roleName),
	})

	if err != nil {
		return role, err
	}

	role.RoleId = *output.Role.RoleId
	role.RoleName = roleName

	instanceRes, err := client.CreateInstanceProfile(context.TODO(), &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(roleName),
		Tags: []types.Tag{
			types.Tag{
				Key:   aws.String(AWS_TAG_KEY),
				Value: aws.String(""),
			},
		},
	})
	if err != nil {
		return role, err
	}

	_, err = client.AddRoleToInstanceProfile(context.TODO(), &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: instanceRes.InstanceProfile.InstanceProfileName,
		RoleName:            aws.String(roleName),
	})
	if err != nil {
		return role, err
	}

	role.Arn = *instanceRes.InstanceProfile.Arn

	return role, nil
}

func (c *AWSClient) GetRoles() ([]types.Role, error) {
	client := iam.NewFromConfig(c.Config)

	output, err := client.ListRoles(context.TODO(), &iam.ListRolesInput{})

	if err != nil {
		return []types.Role{}, err
	}

	return output.Roles, nil
}
