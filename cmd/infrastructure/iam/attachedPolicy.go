package iam

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateAttachedPolicy(ctx *pulumi.Context, name string, policy *iam.Policy, role *iam.Role) (*iam.RolePolicyAttachment, error) {
	attachPolicy, err := iam.NewRolePolicyAttachment(ctx, name, &iam.RolePolicyAttachmentArgs{
		PolicyArn: policy.Arn,
		Role:      role.Name,
	})
	return attachPolicy, err
}
