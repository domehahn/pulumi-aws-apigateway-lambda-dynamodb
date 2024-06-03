package iam

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func AttachedPolicy(ctx *pulumi.Context, name string, role *iam.Role, policyArn pulumi.StringOutput) (*iam.RolePolicyAttachment, error) {
	attachPolicy, err := iam.NewRolePolicyAttachment(ctx, name, &iam.RolePolicyAttachmentArgs{
		PolicyArn: policyArn,
		Role:      role.Name,
	})
	return attachPolicy, err
}
