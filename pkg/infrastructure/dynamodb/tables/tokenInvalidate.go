package tables

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DynamoDbTableTokenInvalidate(ctx *pulumi.Context) (*dynamodb.Table, error) {
	tableTokenInvalidate, err := dynamodb.NewTable(ctx, "TokenInvalidate", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("token"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("invalidate"),
				Type: pulumi.String("N"),
			},
		},
		HashKey: pulumi.String("token"),
		GlobalSecondaryIndexes: dynamodb.TableGlobalSecondaryIndexArray{
			// Define Global Secondary Indexes (GSIs) for other attributes that you want to be able to query against
			// Example GSI (you can add as many as you need):
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String("TokenIndex"),
				HashKey:        pulumi.String("token"),
				ProjectionType: pulumi.String("ALL"),
			},
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String("InvalidateIndex"),
				HashKey:        pulumi.String("invalidate"),
				ProjectionType: pulumi.String("ALL"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
	})
	return tableTokenInvalidate, err
}
