package tables

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DynamoDbTableCartItem(ctx *pulumi.Context) (*dynamodb.Table, error) {
	tableBook, err := dynamodb.NewTable(ctx, "CartItem", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("isbn"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("quantity"),
				Type: pulumi.String("N"),
			},
		},
		HashKey: pulumi.String("isbn"),
		GlobalSecondaryIndexes: dynamodb.TableGlobalSecondaryIndexArray{
			// Define Global Secondary Indexes (GSIs) for other attributes that you want to be able to query against
			// Example GSI (you can add as many as you need):
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String("IsbnIndex"),
				HashKey:        pulumi.String("isbn"),
				ProjectionType: pulumi.String("ALL"),
			},
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String("QuantityIndex"),
				HashKey:        pulumi.String("quantity"),
				ProjectionType: pulumi.String("ALL"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
	})
	return tableBook, err
}
