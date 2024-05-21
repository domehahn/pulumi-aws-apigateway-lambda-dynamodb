package lambda

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateLamdbaFunction(ctx *pulumi.Context, name string, function string, role *iam.Role, table *dynamodb.Table, attachPolicy *iam.RolePolicyAttachment) (*lambda.Function, error) {
	fn, err := lambda.NewFunction(ctx, name, &lambda.FunctionArgs{
		Runtime: pulumi.String("python3.9"),
		Handler: pulumi.String(function),
		Role:    role.Arn,
		Code:    pulumi.NewFileArchive("./function"),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"DYNAMODB_TABLE_NAME": table.Name,
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{attachPolicy}))
	return fn, err
}
