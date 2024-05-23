package tables

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DynamoDbTableBook(ctx *pulumi.Context) (*dynamodb.Table, error) {
	tableBook, err := dynamodb.NewTable(ctx, "Books", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("author"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("title"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("price"),
				Type: pulumi.String("N"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("isbn"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("copiesInStock"),
				Type: pulumi.String("N"),
			},
		},
		HashKey: pulumi.String("isbn"),
		GlobalSecondaryIndexes: dynamodb.TableGlobalSecondaryIndexArray{
			// Define Global Secondary Indexes (GSIs) for other attributes that you want to be able to query against
			// Example GSI (you can add as many as you need):
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String("AuthorIndex"),
				HashKey:        pulumi.String("author"),
				ProjectionType: pulumi.String("ALL"),
			},
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String("TitleIndex"),
				HashKey:        pulumi.String("title"),
				ProjectionType: pulumi.String("ALL"),
			},
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String("PriceIndex"),
				HashKey:        pulumi.String("price"),
				ProjectionType: pulumi.String("ALL"),
			},
			&dynamodb.TableGlobalSecondaryIndexArgs{
				Name:           pulumi.String("CopiesInStockIndex"),
				HashKey:        pulumi.String("copiesInStock"),
				ProjectionType: pulumi.String("ALL"),
			},
		},
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
	})
	return tableBook, err
}
