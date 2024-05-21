package apigateway

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
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
