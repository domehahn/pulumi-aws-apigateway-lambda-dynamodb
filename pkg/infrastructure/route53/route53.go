package route53

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Zone(ctx *pulumi.Context, name string, url string) (*route53.Zone, error) {
	zone, err := route53.NewZone(ctx, name, &route53.ZoneArgs{
		Name: pulumi.String(url),
	})
	return zone, err
}

func Record(ctx *pulumi.Context, name string, zone *route53.Zone, apigateway *apigatewayv2.Api) (*route53.Record, error) {
	route, err := route53.NewRecord(ctx, name, &route53.RecordArgs{
		Name: zone.Name,
		Type: pulumi.String("A"),
		Aliases: route53.RecordAliasArray{
			&route53.RecordAliasArgs{
				Name: apigateway.ApiEndpoint,
				ZoneId: zone.ZoneId.ApplyT(func(zoneId string) string {
					return fmt.Sprintf("%s", zoneId)
				}).(pulumi.StringOutput),
				EvaluateTargetHealth: pulumi.Bool(false),
			},
		},
		ZoneId: pulumi.String("ZLY8HYME6SFDD"), // the ID of your hosted zone in Route53
	})
	return route, err
}
