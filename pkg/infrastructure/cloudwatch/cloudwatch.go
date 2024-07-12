package cloudwatch

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func LogGroup(ctx *pulumi.Context, logGroupName string, path string, name pulumi.StringOutput) (*cloudwatch.LogGroup, error) {
	logGroupPath := name.ApplyT(func(name string) string {
		return path + name
	}).(pulumi.StringOutput)

	logGroup, err := cloudwatch.NewLogGroup(ctx, logGroupName, &cloudwatch.LogGroupArgs{
		Name: logGroupPath.ApplyT(func(name string) string {
			return name
		}).(pulumi.StringOutput),
	})
	return logGroup, err
}

func LogMetrics(ctx *pulumi.Context, metricName string, logGroup *cloudwatch.LogGroup) (*cloudwatch.LogMetricFilter, error) {
	metricFilter, err := cloudwatch.NewLogMetricFilter(ctx, metricName, &cloudwatch.LogMetricFilterArgs{
		LogGroupName: logGroup.Name,
		Pattern:      pulumi.String("[ip, id, user_id, timestamp, request, status_code, size]"),
		MetricTransformation: cloudwatch.LogMetricFilterMetricTransformationArgs{
			Name:      pulumi.String(metricName),
			Namespace: pulumi.String(metricName),
			Value:     pulumi.String("1"),
		},
	})
	return metricFilter, err
}
