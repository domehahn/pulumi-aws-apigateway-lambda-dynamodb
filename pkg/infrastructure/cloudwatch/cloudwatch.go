package cloudwatch

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func LogGroup(ctx *pulumi.Context, name string) (*cloudwatch.LogGroup, error) {
	logGroup, err := cloudwatch.NewLogGroup(ctx, name, nil)
	return logGroup, err
}
