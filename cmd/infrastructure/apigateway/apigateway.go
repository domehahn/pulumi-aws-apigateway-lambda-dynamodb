package apigateway

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateApiGateway(ctx *pulumi.Context, name string) (*apigatewayv2.Api, error) {
	gateway, err := apigatewayv2.NewApi(ctx, name, &apigatewayv2.ApiArgs{
		Name:         pulumi.String(name),
		Description:  pulumi.String("An API Gateway"),
		ProtocolType: pulumi.String("HTTP"),
	})
	return gateway, err
}

func Integration(ctx *pulumi.Context, name string, apigateway *apigatewayv2.Api, lambda *lambda.Function) (*apigatewayv2.Integration, error) {
	integration, err := apigatewayv2.NewIntegration(ctx, name, &apigatewayv2.IntegrationArgs{
		ApiId:                apigateway.ID(),
		IntegrationType:      pulumi.String("AWS_PROXY"),
		IntegrationUri:       lambda.Arn,
		IntegrationMethod:    pulumi.String("POST"),
		PayloadFormatVersion: pulumi.String("2.0"),
	})
	return integration, err
}

func CreateRoute(ctx *pulumi.Context, name string, apigateway *apigatewayv2.Api, routeKey string, integration *apigatewayv2.Integration) (*apigatewayv2.Route, error) {
	route, err := apigatewayv2.NewRoute(ctx, name, &apigatewayv2.RouteArgs{
		ApiId:    apigateway.ID(),
		RouteKey: pulumi.String(fmt.Sprintf("%s", routeKey)),
		Target: integration.ID().ApplyT(func(id pulumi.ID) (string, error) {
			return fmt.Sprintf("integrations/%s", id), nil
		}).(pulumi.StringOutput),
	})
	return route, err
}

func Deploy(ctx *pulumi.Context, name string, apigateway *apigatewayv2.Api, route []*apigatewayv2.Route) (*apigatewayv2.Deployment, error) {
	deployment, err := apigatewayv2.NewDeployment(ctx, name, &apigatewayv2.DeploymentArgs{
		ApiId: apigateway.ID(),
		Triggers: pulumi.StringMap{
			"redeployment": pulumi.String("RedeploymentTrigger"),
		},
	}, pulumi.DependsOn([]pulumi.Resource{route[0], route[1], route[2], route[3], route[4]}))
	return deployment, err
}

func Stage(ctx *pulumi.Context, name string, apiDeployment *apigatewayv2.Deployment, apigateway *apigatewayv2.Api, logGroup *cloudwatch.LogGroup) (*apigatewayv2.Stage, error) {
	stage, err := apigatewayv2.NewStage(ctx, name, &apigatewayv2.StageArgs{
		ApiId:        apigateway.ID(),
		Name:         pulumi.String("dev"),
		DeploymentId: apiDeployment.ID(),
		AccessLogSettings: &apigatewayv2.StageAccessLogSettingsArgs{
			DestinationArn: logGroup.Arn,
			Format:         pulumi.String(`{"requestId":"$context.requestId","ip":"$context.identity.sourceIp","caller":"$context.identity.caller","user":"$context.identity.user","requestTime":"$context.requestTime","httpMethod":"$context.httpMethod","resourcePath":"$context.resourcePath","status":"$context.status","protocol":"$context.protocol","responseLength":"$context.responseLength"}`),
		},
		AutoDeploy: pulumi.Bool(true),
	})
	return stage, err
}
