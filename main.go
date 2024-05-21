package main

import (
	"fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	apigateway2 "pulumi-00/cmd/infrastructure/apigateway"
	"pulumi-00/cmd/infrastructure/cloudwatch"
	"pulumi-00/cmd/infrastructure/dynamodb/tables"
	iam2 "pulumi-00/cmd/infrastructure/iam"
	"pulumi-00/cmd/infrastructure/lambda"
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
		lambdaRole, err := iam2.CreateLambdaRole(ctx, "lambdaRole")
		if err != nil {
			return err
		}

		tableBookPolicy, err := iam2.CreateDynamodbPolicy(ctx, "tableBookPolicy", tableBook)
		if err != nil {
			return err
		}

		tableCartItemPolicy, err := iam2.CreateDynamodbPolicy(ctx, "tableCartItemPolicy", tableCartItem)
		if err != nil {
			return err
		}

		bookAttachedPolicy, err := iam2.CreateAttachedPolicy(ctx, "bookAttachedPolicy", tableBookPolicy, lambdaRole)
		if err != nil {
			return err
		}

		cartItemAttachedPolicy, err := iam2.CreateAttachedPolicy(ctx, "cartItemAttachedPolicy", tableCartItemPolicy, lambdaRole)
		if err != nil {
			return err
		}

		// ### User functions ###
		getBooksFn, err := lambda.CreateLamdbaFunction(ctx, "getBooks", "lambda_book.getBooks", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		// ### Admin functions ###
		createBookFn, err := lambda.CreateLamdbaFunction(ctx, "createBook", "lambda_book.createBook", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		updateBookFn, err := lambda.CreateLamdbaFunction(ctx, "updateBook", "lambda_book.updateBook", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		deleteBookFn, err := lambda.CreateLamdbaFunction(ctx, "deleteBook", "lambda_book.deleteBook", lambdaRole, tableBook, bookAttachedPolicy)
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

		getBooksIntegration, err := apigateway2.Integration(ctx, "getBooksIntegration", apiGateway, getBooksFn)
		if err != nil {
			return err
		}

		createBookIntegration, err := apigateway2.Integration(ctx, "createBookIntegration", apiGateway, createBookFn)
		if err != nil {
			return err
		}

		updateBookIntegration, err := apigateway2.Integration(ctx, "updateBookIntegration", apiGateway, updateBookFn)
		if err != nil {
			return err
		}

		deleteBookIntegration, err := apigateway2.Integration(ctx, "deleteBookIntegration", apiGateway, deleteBookFn)
		if err != nil {
			return err
		}

		_, err = apigateway2.CreateRoute(ctx, "getBooksRoute", apiGateway, "GET /books", getBooksIntegration)

		_, err = apigateway2.CreateRoute(ctx, "createBookRoute", apiGateway, "PUT /createBook", createBookIntegration)

		_, err = apigateway2.CreateRoute(ctx, "updateBookRoute", apiGateway, "PUT /updateBook", updateBookIntegration)

		_, err = apigateway2.CreateRoute(ctx, "deleteBookRoute", apiGateway, "DELETE /deleteBook", deleteBookIntegration)

		apiDeployment, err := apigateway2.Deploy(ctx, "apiDeployment", apiGateway)
		if err != nil {
			return err
		}

		logGroup, err := cloudwatch.CreateLogGroup(ctx, "apigatewayLogGroup")
		if err != nil {
			return err
		}

		stage, err := apigateway2.Stage(ctx, "stage", apiDeployment, apiGateway, logGroup)
		if err != nil {
			return err
		}

		fullApiUrl := apiGateway.ApiEndpoint.ApplyT(func(endpoint string) string {
			return fmt.Sprintf("%s/%s", endpoint, stage.Name)
		}).(pulumi.StringOutput)

		// The URL at which the REST API will be served
		ctx.Export("apiEndpoint", fullApiUrl)
		return nil
	})
}
