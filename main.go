package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	errorhandler "pulumi-00/pkg/error"
	apigateway2 "pulumi-00/pkg/infrastructure/apigateway"
	"pulumi-00/pkg/infrastructure/cloudwatch"
	"pulumi-00/pkg/infrastructure/cognito"
	"pulumi-00/pkg/infrastructure/dynamodb/tables"
	iam2 "pulumi-00/pkg/infrastructure/iam"
	lambdaFn "pulumi-00/pkg/infrastructure/lambda"
)

type RouteConfig struct {
	function    *lambda.Function
	integration *apigatewayv2.Integration
	routePath   string
}

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
		// Lambda role
		lambdaRole, err := iam2.Role(ctx, "lambdaRole", "lambda.amazonaws.com")
		if err != nil {
			return err
		}

		//############################
		// Attached Policies
		bookAttachedPolicy, err := iam2.AttachedPolicy(ctx, "bookLambdaDynamoDbRoleAttachment", lambdaRole, tableBookPolicy.Arn)
		if err != nil {
			return err
		}

		cartItemAttachedPolicy, err := iam2.AttachedPolicy(ctx, "cartItemLambdaDynamoDbRoleAttachment", lambdaRole, tableCartItemPolicy.Arn)
		if err != nil {
			return err
		}

		_, err = iam2.AttachedPolicy(ctx, "lambdaLogGroupRoleAttachment", lambdaRole, pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole").ToStringOutput())
		if err != nil {
			return err
		}

		//############################
		// Lambdas

		// ### User functions ###
		getBooksFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"getBooks",
			"getbooks.getBooks",
			"./function/book",
			lambdaRole,
			tableBook,
			bookAttachedPolicy)

		getBookFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"getBook",
			"getbook.getBook",
			"./function/book",
			lambdaRole, tableBook,
			bookAttachedPolicy)

		getCartItemsFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"cartItems",
			"getcartitems.getCartItems",
			"./function/cartItem",
			lambdaRole,
			tableCartItem,
			cartItemAttachedPolicy)

		addCartItemFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"addCartItem",
			"addcartitem.addCartItem",
			"./function/cartItem",
			lambdaRole,
			tableCartItem,
			cartItemAttachedPolicy)

		deleteCartItemFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"deleteCartItem",
			"deletecartitem.deleteCartItem",
			"./function/cartItem",
			lambdaRole,
			tableCartItem,
			cartItemAttachedPolicy)

		// ### Admin functions ###
		createBookFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"createBook",
			"createbook.createBook",
			"./function/book",
			lambdaRole,
			tableBook,
			bookAttachedPolicy)

		updateBookFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"updateBook",
			"updatebook.updateBook",
			"./function/book",
			lambdaRole,
			tableBook,
			bookAttachedPolicy)

		deleteBookFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"deleteBook",
			"deletebook.deleteBook",
			"./function/book",
			lambdaRole,
			tableBook,
			bookAttachedPolicy)

		// ### Authorization and Authentication functions ###
		loginFn, err := lambdaFn.LamdbaFunctionShort(
			ctx,
			"login",
			"login.login",
			"./function/auth",
			lambdaRole)

		logoutFn, err := lambdaFn.LamdbaFunctionShort(
			ctx,
			"logout",
			"logout.logout",
			"./function/auth",
			lambdaRole)

		if err != nil {
			errorhandler.HandlingError("Error creating lambda.")
		}

		//############################
		// Route53 zone
		/*route53Zone, err := route53.Zone(ctx, "Zone", "aws-lerngruppe-2.com")
		if err != nil {
			return err
		}*/

		//############################
		// Route53 record
		/*_, err = route53.Record(ctx, "Record", route53Zone, apiGateway)
		if err != nil {
			return err
		}*/

		//############################
		// Cognito Userpool
		userPool, err := cognito.UserPool(ctx, "UserPool")
		if err != nil {
			return err
		}

		//############################
		// Cognito Userpool
		userPoolClient, err := cognito.UserPoolClient(ctx, "UserPoolClient", userPool)
		if err != nil {
			return err
		}

		//############################
		// Api LogGroup
		logGroup, err := cloudwatch.LogGroup(ctx, "apigatewayLogGroup", "/aws/apigateway/log-group")
		if err != nil {
			return err
		}

		//############################
		// Api Gateway role
		gatewayRole, err := iam2.Role(ctx, "apigatewayRole", "apigateway.amazonaws.com")
		if err != nil {
			errorhandler.HandlingError("Error creating API-Gateway role.")
		}

		//############################
		// Api Gateway role Policies
		_, err = iam2.AttachedPolicy(ctx, "apigatewayLogGroupRoleAttachment", gatewayRole, pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs").ToStringOutput())

		//############################
		// Api Gateway
		apiGateway, err := apigateway2.ApiGateway(ctx, "ApiGateway")
		if err != nil {
			errorhandler.HandlingError("Error creating API-Gateway.")
		}

		//############################
		// Api Gateway Lambda Permission
		_, err = lambdaFn.LambdaPermission(ctx, "bookPermission", getBookFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "booksPermission", getBooksFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "createBookPermission", createBookFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "updateBookPermission", updateBookFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "deleteBookPermission", deleteBookFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "cartItemsPermission", getCartItemsFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "addCartItemPermission", addCartItemFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "deleteCartItemPermission", deleteCartItemFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "loginPermission", loginFn, apiGateway)

		_, err = lambdaFn.LambdaPermission(ctx, "logoutPermission", logoutFn, apiGateway)

		//############################
		// Api Authorizer
		authorizer, err := apigateway2.Authorizer(ctx, "Authorizer", apiGateway, userPool, userPoolClient)
		if err != nil {
			errorhandler.HandlingError("Error creating authorizer.")
		}

		//############################
		// Api Integrations
		bookIntegration, err := apigateway2.Integration(ctx, "bookIntegration", apiGateway, getBookFn)

		booksIntegration, err := apigateway2.Integration(ctx, "booksIntegration", apiGateway, getBooksFn)

		createBookIntegration, err := apigateway2.Integration(ctx, "createBookIntegration", apiGateway, createBookFn)

		updateBookIntegration, err := apigateway2.Integration(ctx, "updateBookIntegration", apiGateway, updateBookFn)

		deleteBookIntegration, err := apigateway2.Integration(ctx, "deleteBookIntegration", apiGateway, deleteBookFn)

		cartItemsIntegration, err := apigateway2.Integration(ctx, "cartItemsIntegration", apiGateway, getCartItemsFn)

		addCartItemIntegration, err := apigateway2.Integration(ctx, "addCartItemIntegration", apiGateway, addCartItemFn)

		deleteCartItemIntegration, err := apigateway2.Integration(ctx, "deleteCartItemIntegration", apiGateway, deleteCartItemFn)

		loginIntegration, err := apigateway2.Integration(ctx, "loginIntegration", apiGateway, loginFn)

		logoutIntegration, err := apigateway2.Integration(ctx, "logoutIntegration", apiGateway, logoutFn)
		if err != nil {
			errorhandler.HandlingError("Error creating integration.")
		}

		bookRoute, err := apigateway2.Route(ctx, "bookRoute", apiGateway, "GET /books/{isbn}", bookIntegration, authorizer)

		booksRoute, err := apigateway2.Route(ctx, "booksRoute", apiGateway, "GET /books", booksIntegration, authorizer)

		createBookRoute, err := apigateway2.Route(ctx, "createBookRoute", apiGateway, "POST /books", createBookIntegration, authorizer)

		updateBookRoute, err := apigateway2.Route(ctx, "updateBookRoute", apiGateway, "PATCH /books/{isbn}", updateBookIntegration, authorizer)

		deleteBookRoute, err := apigateway2.Route(ctx, "deleteBookRoute", apiGateway, "DELETE /books/{isbn}", deleteBookIntegration, authorizer)

		cartItemsRoute, err := apigateway2.Route(ctx, "cartItemsRoute", apiGateway, "GET /cart", cartItemsIntegration, authorizer)

		addCartItemRoute, err := apigateway2.Route(ctx, "addCartItemRoute", apiGateway, "POST /cart", addCartItemIntegration, authorizer)

		deleteCartItemRoute, err := apigateway2.Route(ctx, "deleteCartItemRoute", apiGateway, "DELETE /cart/{isbn}", deleteCartItemIntegration, authorizer)

		loginRoute, err := apigateway2.RouteWithoutAuthorizer(ctx, "loginRoute", apiGateway, "POST /users/login", loginIntegration)

		logoutRoute, err := apigateway2.Route(ctx, "logoutRoute", apiGateway, "POST /users/logout", logoutIntegration, authorizer)
		if err != nil {
			errorhandler.HandlingError("Error creating route.")
		}

		//############################
		// Api Deployment
		apiDeployment, err := apigateway2.Deploy(ctx, "apiDeployment", apiGateway, []*apigatewayv2.Route{
			booksRoute,
			bookRoute,
			createBookRoute,
			updateBookRoute,
			deleteBookRoute,
			cartItemsRoute,
			addCartItemRoute,
			deleteCartItemRoute,
			loginRoute,
			logoutRoute})
		if err != nil {
			errorhandler.HandlingError("Error creating deployment.")
		}

		//############################
		// Api Stage
		stage, err := apigateway2.Stage(ctx, "stage", apiDeployment, apiGateway, logGroup)
		if err != nil {
			errorhandler.HandlingError("Error creating stage.")
		}

		fullApiUrl := apiGateway.ApiEndpoint.ApplyT(func(endpoint string) string {
			return fmt.Sprintf("%s/%s", endpoint, stage.Name)
		}).(pulumi.StringOutput)

		// The URL at which the REST API will be served
		ctx.Export("apiEndpoint", fullApiUrl)

		return nil
	})
}
