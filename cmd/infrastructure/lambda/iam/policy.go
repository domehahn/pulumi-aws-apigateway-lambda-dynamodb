package iam

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateDynamodbPolicy(ctx *pulumi.Context, name string, table *dynamodb.Table) (*iam.Policy, error) {
	// Create an IAM Policy for Lambda that grants access to the DynamoDB table.
	policy, err := iam.NewPolicy(ctx, name, &iam.PolicyArgs{
		Description: pulumi.String("IAM policy for Lambda to access DynamoDB table"),
		Policy: table.Arn.ApplyT(func(arn string) (string, error) {
			// Construct the policy document with the correct ARN for the DynamoDB table.
			policyDocument := `{
					"Version": "2012-10-17",
					"Statement": [{
						"Sid": "ReadWriteTable",
						"Effect": "Allow",
						"Action": [
				"dynamodb:BatchGetItem",
                "dynamodb:GetItem",
                "dynamodb:Query",
                "dynamodb:Scan",
                "dynamodb:BatchWriteItem",
                "dynamodb:PutItem",
                "dynamodb:UpdateItem"
				],
						"Resource": "` + arn + `"
					}]
				}`
			return policyDocument, nil
		}).(pulumi.StringOutput),
	})
	return policy, err
}
