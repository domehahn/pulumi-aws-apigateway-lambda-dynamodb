package iam

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func RolePolicy(ctx *pulumi.Context, name string, role *iam.Role) (*iam.RolePolicy, error) {
	logPolicy, err := iam.NewRolePolicy(ctx, name, &iam.RolePolicyArgs{
		Role: role.Name,
		Policy: pulumi.String(`{
                "Version": "2012-10-17",
                "Statement": [{
                    "Effect": "Allow",
                    "Action": [
                        "logs:CreateLogGroup",
                        "logs:CreateLogStream",
                        "logs:PutLogEvents"
                    ],
                    "Resource": "arn:aws:logs:*:*:*"
                }]
            }`),
	})
	return logPolicy, err
}
