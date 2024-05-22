package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	apigateway2 "pulumi-00/cmd/infrastructure/apigateway"
	"pulumi-00/cmd/infrastructure/cloudwatch"
	"pulumi-00/cmd/infrastructure/dynamodb/tables"
	iam2 "pulumi-00/cmd/infrastructure/iam"
	"pulumi-00/cmd/infrastructure/lambda"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		//############################
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

		//############################
		// Lambdarole
		lambdaRole, err := iam2.CreateLambdaRole(ctx, "lambdaRole")
		if err != nil {
			return err
		}

		//############################
		// Attached Policies
		bookAttachedPolicy, err := iam2.CreateAttachedPolicy(ctx, "bookAttachedPolicy", tableBookPolicy, lambdaRole)
		if err != nil {
			return err
		}

		cartItemAttachedPolicy, err := iam2.CreateAttachedPolicy(ctx, "cartItemAttachedPolicy", tableCartItemPolicy, lambdaRole)
		if err != nil {
			return err
		}

		//############################
		// Lambdas

		// ### User functions ###
		getBooksFn, err := lambda.CreateLamdbaFunction(ctx, "getBooks", "getbooks.getBooks", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		getBookFn, err := lambda.CreateLamdbaFunction(ctx, "getBook", "getbook.getBook", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		getCartItemsFn, err := lambda.CreateLamdbaFunction(ctx, "cartItems", "getcartitems.getCartItems", "./function/cartItem", lambdaRole, tableCartItem, cartItemAttachedPolicy)
		if err != nil {
			return err
		}

		addCartItemFn, err := lambda.CreateLamdbaFunction(ctx, "addCartItem", "addcartitem.addCartItem", "./function/cartItem", lambdaRole, tableCartItem, cartItemAttachedPolicy)
		if err != nil {
			return err
		}

		deleteCartItemFn, err := lambda.CreateLamdbaFunction(ctx, "deleteCartItem", "deletecartitem.deleteCartItem", "./function/cartItem", lambdaRole, tableCartItem, cartItemAttachedPolicy)
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

		//############################
		// Api Gateway
		apiGateway, err := apigateway2.CreateApiGateway(ctx, "ApiGateway")
		if err != nil {
			return err
		}

		//############################
		// Api Integrations
		getBooksIntegration, err := apigateway2.Integration(ctx, "getBooksIntegration", apiGateway, getBooksFn)
		if err != nil {
			return err
		}

		getBookIntegration, err := apigateway2.Integration(ctx, "getBookIntegration", apiGateway, getBookFn)
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

		getCartItemsIntegration, err := apigateway2.Integration(ctx, "getCartItemsIntegration", apiGateway, getCartItemsFn)
		if err != nil {
			return err
		}

		addCartItemIntegration, err := apigateway2.Integration(ctx, "addCartItemIntegration", apiGateway, addCartItemFn)
		if err != nil {
			return err
		}

		deleteCartItemIntegration, err := apigateway2.Integration(ctx, "deleteCartItemIntegration", apiGateway, deleteCartItemFn)
		if err != nil {
			return err
		}

		//############################
		// Api Routes
		booksRoute, err := apigateway2.CreateRoute(ctx, "getBooksRoute", apiGateway, "GET /books", getBooksIntegration)
		if err != nil {
			return err
		}

		bookRoute, err := apigateway2.CreateRoute(ctx, "getBookRoute", apiGateway, "GET /books/{isbn}", getBookIntegration)
		if err != nil {
			return err
		}

		createBookRoute, err := apigateway2.CreateRoute(ctx, "createBookRoute", apiGateway, "POST /books", createBookIntegration)
		if err != nil {
			return err
		}

		updateBookRoute, err := apigateway2.CreateRoute(ctx, "updateBookRoute", apiGateway, "PATCH /books/{isbn}", updateBookIntegration)
		if err != nil {
			return err
		}

		deleteBookRoute, err := apigateway2.CreateRoute(ctx, "deleteBookRoute", apiGateway, "DELETE /books/{isbn}", deleteBookIntegration)
		if err != nil {
			return err
		}

		getCartItemsRoute, err := apigateway2.CreateRoute(ctx, "getCartItemsRoute", apiGateway, "GET /cart", getCartItemsIntegration)
		if err != nil {
			return err
		}

		addCartItemRoute, err := apigateway2.CreateRoute(ctx, "addCartItemRoute", apiGateway, "POST /cart", addCartItemIntegration)
		if err != nil {
			return err
		}

		deleteCartItemRoute, err := apigateway2.CreateRoute(ctx, "deleteCartItemRoute", apiGateway, "DELETE /cart/{isbn}", deleteCartItemIntegration)
		if err != nil {
			return err
		}

		//############################
		// Api Deployment
		apiDeployment, err := apigateway2.Deploy(ctx, "apiDeployment", apiGateway, []*apigatewayv2.Route{booksRoute, bookRoute, createBookRoute, updateBookRoute, deleteBookRoute, getCartItemsRoute, addCartItemRoute, deleteCartItemRoute})
		if err != nil {
			return err
		}

		//############################
		// Api LogGroup
		logGroup, err := cloudwatch.CreateLogGroup(ctx, "apigatewayLogGroup")
		if err != nil {
			return err
		}

		//############################
		// Api Stage
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
