package lambda

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func LamdbaFunction(ctx *pulumi.Context, name string, function string, path string, role *iam.Role, table *dynamodb.Table, environmentVariables pulumi.StringMap, attachPolicies []*iam.RolePolicyAttachment) (*lambda.Function, error) {
	fn, err := lambda.NewFunction(ctx, name, &lambda.FunctionArgs{
		Runtime: pulumi.String("python3.9"),
		Handler: pulumi.String(function),
		Role:    role.Arn,
		Code:    pulumi.NewFileArchive(path),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: environmentVariables,
		},
	}, pulumi.DependsOn(func() []pulumi.Resource {
		deps := make([]pulumi.Resource, len(attachPolicies))
		for i, policy := range attachPolicies {
			deps[i] = policy
		}
		return deps
	}()))
	return fn, err
}

func LamdbaFunctionShort(ctx *pulumi.Context, name string, function string, path string, role *iam.Role) (*lambda.Function, error) {
	fn, err := lambda.NewFunction(ctx, name, &lambda.FunctionArgs{
		Runtime: pulumi.String("python3.9"),
		Handler: pulumi.String(function),
		Role:    role.Arn,
		Code:    pulumi.NewFileArchive(path),
	})
	return fn, err
}

func LambdaPermission(ctx *pulumi.Context, name string, lambdaFn *lambda.Function, apigateway *apigatewayv2.Api) (*lambda.Permission, error) {
	permission, err := lambda.NewPermission(ctx, name, &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  lambdaFn.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("%s/*/*", apigateway.ExecutionArn),
	})
	lambda.NewInvocation(ctx, name, &lambda.InvocationArgs{})
	return permission, err
}
