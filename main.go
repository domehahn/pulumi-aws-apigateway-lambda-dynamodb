package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	apigateway2 "pulumi-00/cmd/infrastructure/apigateway"
	"pulumi-00/cmd/infrastructure/dynamodb/tables"
	"pulumi-00/cmd/infrastructure/lambda"
	"pulumi-00/cmd/infrastructure/lambda/iam"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		//Create infrastructure
		tableBook, err := tables.CreateDynamoDBTableBook(ctx)
		if err != nil {
			return err
		}

		tableCartItem, err := tables.CreateDynamoDBTableCartItem(ctx)
		if err != nil {
			return err
		}

		//Create lambda  role
		lambdaRole, err := iam.CreateLambdaRole(ctx, "lambdaRole")
		if err != nil {
			return err
		}

		tableBookPolicy, err := iam.CreateDynamodbPolicy(ctx, "tableBookPolicy", tableBook)
		if err != nil {
			return err
		}

		tableCartItemPolicy, err := iam.CreateDynamodbPolicy(ctx, "tableCartItemPolicy", tableCartItem)
		if err != nil {
			return err
		}

		bookAttachedPolicy, err := iam.CreateAttachedPolicy(ctx, "bookAttachedPolicy", tableBookPolicy, lambdaRole)
		if err != nil {
			return err
		}

		cartItemAttachedPolicy, err := iam.CreateAttachedPolicy(ctx, "cartItemAttachedPolicy", tableCartItemPolicy, lambdaRole)
		if err != nil {
			return err
		}

		_, err = lambda.CreateLamdbaFunction(ctx, "getBooks", "lambda_book.getBooks", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		_, err = lambda.CreateLamdbaFunction(ctx, "createBook", "lambda_book.createBook", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		_, err = lambda.CreateLamdbaFunction(ctx, "cartItems", "lambda_cartItem.getCartItems", lambdaRole, tableCartItem, cartItemAttachedPolicy)
		if err != nil {
			return err
		}

		apiGateway, err := apigateway2.CreateApiGateway(ctx, "ApiGateway")
		if err != nil {
			return err
		}

		// The URL at which the REST API will be served
		ctx.Export("apiUrl", apiGateway.ApiEndpoint)
		return nil
	})
}
