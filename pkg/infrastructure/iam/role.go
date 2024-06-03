package iam

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Role(ctx *pulumi.Context, name string, serviceName string) (*iam.Role, error) {
	// IAM role for the Lambda function
	lambdaRole, err := iam.NewRole(ctx, name, &iam.RoleArgs{
		AssumeRolePolicy: pulumi.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
						"Service": "%s"
					}
				}]
			}`, serviceName),
	})
	return lambdaRole, err
}
