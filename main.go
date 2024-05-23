package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	apigateway2 "pulumi-00/pkg/infrastructure/apigateway"
	"pulumi-00/pkg/infrastructure/cloudwatch"
	"pulumi-00/pkg/infrastructure/dynamodb/tables"
	iam2 "pulumi-00/pkg/infrastructure/iam"
	"pulumi-00/pkg/infrastructure/lambda"
	"pulumi-00/pkg/infrastructure/route53"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		//############################
		// DynamoDb
		tableBook, err := tables.DynamoDbTableBook(ctx)
		if err != nil {
			return err
		}

		tableBookPolicy, err := iam2.DynamoDbPolicy(ctx, "tableBookPolicy", tableBook)
		if err != nil {
			return err
		}

		tableCartItem, err := tables.DynamoDbTableCartItem(ctx)
		if err != nil {
			return err
		}

		tableCartItemPolicy, err := iam2.DynamoDbPolicy(ctx, "tableCartItemPolicy", tableCartItem)
		if err != nil {
			return err
		}

		//############################
		// Role
		lambdaRole, err := iam2.LambdaRole(ctx, "lambdaRole")
		if err != nil {
			return err
		}

		//############################
		// Role Policies
		_, err = iam2.RolePolicy(ctx, "lambdaLogPolicy", lambdaRole)

		//############################
		// Attached Policies
		bookAttachedPolicy, err := iam2.AttachedPolicy(ctx, "bookAttachedPolicy", tableBookPolicy, lambdaRole)
		if err != nil {
			return err
		}

		cartItemAttachedPolicy, err := iam2.AttachedPolicy(ctx, "cartItemAttachedPolicy", tableCartItemPolicy, lambdaRole)
		if err != nil {
			return err
		}

		//############################
		// Lambdas

		// ### User functions ###
		getBooksFn, err := lambda.LamdbaFunction(ctx, "getBooks", "getbooks.getBooks", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		getBookFn, err := lambda.LamdbaFunction(ctx, "getBook", "getbook.getBook", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		getCartItemsFn, err := lambda.LamdbaFunction(ctx, "cartItems", "getcartitems.getCartItems", "./function/cartItem", lambdaRole, tableCartItem, cartItemAttachedPolicy)
		if err != nil {
			return err
		}

		addCartItemFn, err := lambda.LamdbaFunction(ctx, "addCartItem", "addcartitem.addCartItem", "./function/cartItem", lambdaRole, tableCartItem, cartItemAttachedPolicy)
		if err != nil {
			return err
		}

		deleteCartItemFn, err := lambda.LamdbaFunction(ctx, "deleteCartItem", "deletecartitem.deleteCartItem", "./function/cartItem", lambdaRole, tableCartItem, cartItemAttachedPolicy)
		if err != nil {
			return err
		}

		// ### Admin functions ###
		createBookFn, err := lambda.LamdbaFunction(ctx, "createBook", "createbook.createBook", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		updateBookFn, err := lambda.LamdbaFunction(ctx, "updateBook", "updatebook.updateBook", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		deleteBookFn, err := lambda.LamdbaFunction(ctx, "deleteBook", "deletebook.deleteBook", "./function/book", lambdaRole, tableBook, bookAttachedPolicy)
		if err != nil {
			return err
		}

		//############################
		// Api Gateway
		apiGateway, err := apigateway2.ApiGateway(ctx, "ApiGateway")
		if err != nil {
			return err
		}

		//############################
		// Route53 zone
		route53Zone, err := route53.Zone(ctx, "Zone", "aws-lerngruppe-2.com")
		if err != nil {
			return err
		}

		//############################
		// Route53 record
		_, err = route53.Record(ctx, "Record", route53Zone, apiGateway)
		if err != nil {
			return err
		}

		//############################
		// Api Gateway Lambda Permission
		_, err = lambda.LambdaPermission(ctx, "bookPermission", getBookFn, apiGateway)

		_, err = lambda.LambdaPermission(ctx, "booksPermission", getBooksFn, apiGateway)

		_, err = lambda.LambdaPermission(ctx, "createBookPermission", createBookFn, apiGateway)

		_, err = lambda.LambdaPermission(ctx, "updateBookPermission", updateBookFn, apiGateway)

		_, err = lambda.LambdaPermission(ctx, "deleteBookPermission", deleteBookFn, apiGateway)

		_, err = lambda.LambdaPermission(ctx, "getCartItemsPermission", getCartItemsFn, apiGateway)

		_, err = lambda.LambdaPermission(ctx, "addCartItemsPermission", addCartItemFn, apiGateway)

		_, err = lambda.LambdaPermission(ctx, "deleteCartItemPermission", deleteCartItemFn, apiGateway)

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
		booksRoute, err := apigateway2.Route(ctx, "getBooksRoute", apiGateway, "GET /books", getBooksIntegration)
		if err != nil {
			return err
		}

		bookRoute, err := apigateway2.Route(ctx, "getBookRoute", apiGateway, "GET /books/{isbn}", getBookIntegration)
		if err != nil {
			return err
		}

		createBookRoute, err := apigateway2.Route(ctx, "createBookRoute", apiGateway, "POST /books", createBookIntegration)
		if err != nil {
			return err
		}

		updateBookRoute, err := apigateway2.Route(ctx, "updateBookRoute", apiGateway, "PATCH /books/{isbn}", updateBookIntegration)
		if err != nil {
			return err
		}

		deleteBookRoute, err := apigateway2.Route(ctx, "deleteBookRoute", apiGateway, "DELETE /books/{isbn}", deleteBookIntegration)
		if err != nil {
			return err
		}

		getCartItemsRoute, err := apigateway2.Route(ctx, "getCartItemsRoute", apiGateway, "GET /cart", getCartItemsIntegration)
		if err != nil {
			return err
		}

		addCartItemRoute, err := apigateway2.Route(ctx, "addCartItemRoute", apiGateway, "POST /cart", addCartItemIntegration)
		if err != nil {
			return err
		}

		deleteCartItemRoute, err := apigateway2.Route(ctx, "deleteCartItemRoute", apiGateway, "DELETE /cart/{isbn}", deleteCartItemIntegration)
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
		logGroup, err := cloudwatch.LogGroup(ctx, "apigatewayLogGroup")
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
