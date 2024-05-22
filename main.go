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
		// DynamoDb
		tableBook, err := tables.CreateDynamoDBTableBook(ctx)
		if err != nil {
			return err
		}

		tableBookPolicy, err := iam2.CreateDynamodbPolicy(ctx, "tableBookPolicy", tableBook)
		if err != nil {
			return err
		}

		tableCartItem, err := tables.CreateDynamoDBTableCartItem(ctx)
		if err != nil {
			return err
		}

		tableCartItemPolicy, err := iam2.CreateDynamodbPolicy(ctx, "tableCartItemPolicy", tableCartItem)
		if err != nil {
			return err
		}

		// Lambdarole
		lambdaRole, err := iam2.CreateLambdaRole(ctx, "lambdaRole")
		if err != nil {
			return err
		}

		// Attached Policies
		bookAttachedPolicy, err := iam2.CreateAttachedPolicy(ctx, "bookAttachedPolicy", tableBookPolicy, lambdaRole)
		if err != nil {
			return err
		}

		cartItemAttachedPolicy, err := iam2.CreateAttachedPolicy(ctx, "cartItemAttachedPolicy", tableCartItemPolicy, lambdaRole)
		if err != nil {
			return err
		}

		// Lambdas

		// ### User functions ###
		getBooksFn, err := lambda.CreateLamdbaFunction(ctx, "getBooks", "getbook.getBooks", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		// ### Admin functions ###
		createBookFn, err := lambda.CreateLamdbaFunction(ctx, "createBook", "createbook.createBook", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		updateBookFn, err := lambda.CreateLamdbaFunction(ctx, "updateBook", "updatebook.updateBook", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		deleteBookFn, err := lambda.CreateLamdbaFunction(ctx, "deleteBook", "deletebook.deleteBook", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		_, err = lambda.CreateLamdbaFunction(ctx, "cartItems", "getcart.getCartItems", "./function/cart", lambdaRole, tableCartItem, cartItemAttachedPolicy)
		if err != nil {
			return err
		}

		// Api Gateway
		apiGateway, err := apigateway2.CreateApiGateway(ctx, "ApiGateway")
		if err != nil {
			return err
		}

		// Api Integrations
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

		// Api Routes
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
