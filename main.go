package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"pulumi-00/pkg/environment"
	errorhandler "pulumi-00/pkg/error"
	apigateway2 "pulumi-00/pkg/infrastructure/apigateway"
	"pulumi-00/pkg/infrastructure/cloudwatch"
	"pulumi-00/pkg/infrastructure/cognito"
	"pulumi-00/pkg/infrastructure/dynamodb/tables"
	iam2 "pulumi-00/pkg/infrastructure/iam"
	lambdaFn "pulumi-00/pkg/infrastructure/lambda"
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

		tableTokenInvalidate, err := tables.DynamoDbTableTokenInvalidate(ctx)
		if err != nil {
			return err
		}

		tableTokenInvalidatePolicy, err := iam2.DynamoDbPolicy(ctx, "tableTokenInvalidate", tableTokenInvalidate)
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
		bookAttachedPolicy, err := iam2.AttachedPolicy(
			ctx,
			"bookLambdaDynamoDbRoleAttachment",
			lambdaRole,
			tableBookPolicy.Arn)
		if err != nil {
			return err
		}

		cartItemAttachedPolicy, err := iam2.AttachedPolicy(
			ctx,
			"cartItemLambdaDynamoDbRoleAttachment",
			lambdaRole,
			tableCartItemPolicy.Arn)
		if err != nil {
			return err
		}

		tokenInvalidateAttachedPolicy, err := iam2.AttachedPolicy(
			ctx,
			"tokenInvalidateLambdaDynamoDbRoleAttachment",
			lambdaRole,
			tableTokenInvalidatePolicy.Arn)
		if err != nil {
			return err
		}

		_, err = iam2.AttachedPolicy(
			ctx,
			"lambdaRoleAttachment",
			lambdaRole,
			pulumi.String("arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole").ToStringOutput())
		if err != nil {
			return err
		}

		//############################
		// Lambdas

		// ### User functions ###
		cartItemTableEnv := pulumi.StringMap{
			"DYNAMODB_TABLE_NAME": tableCartItem.Name,
		}

		bookTableEnv := pulumi.StringMap{
			"DYNAMODB_TABLE_NAME": tableBook.Name,
		}

		tokenInvalidateTableEnv := pulumi.StringMap{
			"DYNAMODB_TABLE_NAME": tableTokenInvalidate.Name,
		}

		getBooksFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"getBooks",
			"getbooks.getBooks",
			"./function/book",
			lambdaRole,
			tableBook,
			bookTableEnv,
			[]*iam.RolePolicyAttachment{bookAttachedPolicy, tokenInvalidateAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "getBooksLogGroup", "/aws/lambda/", getBooksFn.Name)
		if err != nil {
			return err
		}

		getBookFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"getBook",
			"getbook.getBook",
			"./function/book",
			lambdaRole, tableBook, bookTableEnv,
			[]*iam.RolePolicyAttachment{bookAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "getBookLogGroup", "/aws/lambda/", getBookFn.Name)
		if err != nil {
			return err
		}

		getCartItemsFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"cartItems",
			"getcartitems.getCartItems",
			"./function/cartItem",
			lambdaRole,
			tableCartItem, cartItemTableEnv,
			[]*iam.RolePolicyAttachment{cartItemAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "cartItemsLogGroup", "/aws/lambda/", getCartItemsFn.Name)
		if err != nil {
			return err
		}

		multiTableEnv := pulumi.StringMap{
			"DYNAMODB_TABLE_NAME": tableCartItem.Name,
			"UPDATE_TABLE_NAME":   tableBook.Name,
		}
		addCartItemFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"addCartItem",
			"addcartitem.addCartItem",
			"./function/cartItem",
			lambdaRole,
			tableCartItem, multiTableEnv,
			[]*iam.RolePolicyAttachment{cartItemAttachedPolicy, bookAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "addCartItemLogGroup", "/aws/lambda/", addCartItemFn.Name)
		if err != nil {
			return err
		}

		deleteCartItemFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"deleteCartItem",
			"deletecartitem.deleteCartItem",
			"./function/cartItem",
			lambdaRole,
			tableCartItem, multiTableEnv,
			[]*iam.RolePolicyAttachment{cartItemAttachedPolicy, bookAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "deleteCartItemLogGroup", "/aws/lambda/", deleteCartItemFn.Name)
		if err != nil {
			return err
		}

		// ### Admin functions ###
		createBookFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"createBook",
			"createbook.createBook",
			"./function/book",
			lambdaRole,
			tableBook, bookTableEnv,
			[]*iam.RolePolicyAttachment{bookAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "createBookLogGroup", "/aws/lambda/", createBookFn.Name)
		if err != nil {
			return err
		}

		updateBookFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"updateBook",
			"updatebook.updateBook",
			"./function/book",
			lambdaRole,
			tableBook, bookTableEnv,
			[]*iam.RolePolicyAttachment{bookAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "updateBookLogGroup", "/aws/lambda/", updateBookFn.Name)
		if err != nil {
			return err
		}

		multiTableEnv = pulumi.StringMap{
			"DYNAMODB_TABLE_NAME": tableBook.Name,
			"UPDATE_TABLE_NAME":   tableCartItem.Name,
		}
		deleteBookFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"deleteBook",
			"deletebook.deleteBook",
			"./function/book",
			lambdaRole,
			tableBook, multiTableEnv,
			[]*iam.RolePolicyAttachment{bookAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "deleteBookLogGroup", "/aws/lambda/", deleteBookFn.Name)
		if err != nil {
			return err
		}

		// ### Authorization and Authentication functions ###
		loginFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"login",
			"login.login",
			"./function/auth",
			lambdaRole,
			tableTokenInvalidate,
			tokenInvalidateTableEnv,
			[]*iam.RolePolicyAttachment{tokenInvalidateAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "loginLogGroup", "/aws/lambda/", loginFn.Name)
		if err != nil {
			return err
		}

		logoutFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"logout",
			"logout.logout",
			"./function/auth",
			lambdaRole,
			tableTokenInvalidate,
			tokenInvalidateTableEnv,
			[]*iam.RolePolicyAttachment{tokenInvalidateAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "logoutLogGroup", "/aws/lambda/", logoutFn.Name)
		if err != nil {
			return err
		}

		authorizeFn, err := lambdaFn.LamdbaFunction(
			ctx,
			"authorize",
			"authorize.authorize",
			"./function/auth",
			lambdaRole,
			tableTokenInvalidate,
			tokenInvalidateTableEnv,
			[]*iam.RolePolicyAttachment{tokenInvalidateAttachedPolicy})
		if err != nil {
			return err
		}

		_, err = cloudwatch.LogGroup(ctx, "authorizeLogGroup", "/aws/lambda/", authorizeFn.Name)
		if err != nil {
			return err
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
		// Cognito User
		username, err := environment.ViperGetEnvVariable("sensitive.user.username")
		if err != nil {
			return err
		}

		password, err := environment.ViperGetEnvVariable("sensitive.user.password")
		if err != nil {
			return err
		}

		email, err := environment.ViperGetEnvVariable("sensitive.user.email")
		if err != nil {
			return err
		}

		sensitiveUser := cognito.SensitiveUser{
			Username: username,
			Password: password,
			Email:    email,
		}
		_, err = cognito.User(ctx, "User", sensitiveUser, userPool)
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
		// Api Gateway
		apiGateway, err := apigateway2.ApiGateway(ctx, "ApiGateway")
		if err != nil {
			errorhandler.HandlingError("Error creating API-Gateway.")
		}

		//############################
		// Api LogGroup
		logGroup, err := cloudwatch.LogGroup(ctx, "apiGatewayLogGroup", "/aws/apigateway/", apiGateway.Name)
		if err != nil {
			return err
		}

		//############################
		// Api Gateway role
		gatewayRole, err := iam2.Role(ctx, "apiGatewayRole", "apigateway.amazonaws.com")
		if err != nil {
			errorhandler.HandlingError("Error creating API-Gateway role.")
		}

		//############################
		// Api Gateway role Policies
		_, err = iam.NewRolePolicyAttachment(ctx, "apiGatewayLogRolePolicy", &iam.RolePolicyAttachmentArgs{
			Role:      gatewayRole.Name,
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/CloudWatchLogsFullAccess"),
		})
		if err != nil {
			return err
		}

		//############################
		// Api Gateway Lambda Permission

		lambdaPermissions := []struct {
			name     string
			lambdaFn *lambda.Function
		}{
			{"bookPermission", getBookFn},
			{"booksPermission", getBooksFn},
			{"createBookPermission", createBookFn},
			{"updateBookPermission", updateBookFn},
			{"deleteBookPermission", deleteBookFn},
			{"cartItemsPermission", getCartItemsFn},
			{"addCartItemPermission", addCartItemFn},
			{"deleteCartItemPermission", deleteCartItemFn},
			{"loginPermission", loginFn},
			{"logoutPermission", logoutFn},
			{"authorizePermission", authorizeFn},
		}

		for _, lambdaPermission := range lambdaPermissions {
			_, err = lambdaFn.LambdaPermission(ctx, lambdaPermission.name, lambdaPermission.lambdaFn, apiGateway)
			if err != nil {
				return err
			}
		}

		//############################
		// Api Authorizer
		lambdaAuthorizer, err := apigateway2.LambdaAuthorizer(ctx, "AuthorizerTest", apiGateway, authorizeFn)
		if err != nil {
			errorhandler.HandlingError("Error creating authorizer.")
		}

		//############################
		// Api Integrations and Routes
		var routes []*apigatewayv2.Route

		integrations := []struct {
			name       string
			httpMethod string
			path       string
			function   *lambda.Function
		}{
			{"bookIntegration", "GET", "/books/{isbn}", getBookFn},
			{"booksIntegration", "GET", "/books", getBooksFn},
			{"createBookIntegration", "POST", "/books", createBookFn},
			{"updateBookIntegration", "PATCH", "/books/{isbn}", updateBookFn},
			{"deleteBookIntegration", "DELETE", "/books/{isbn}", deleteBookFn},
			{"cartItemsIntegration", "GET", "/cart", getCartItemsFn},
			{"addCartItemIntegration", "POST", "/cart", addCartItemFn},
			{"deleteCartItemIntegration", "DELETE", "/cart/{isbn}", deleteCartItemFn},
			{"loginIntegration", "POST", "/users/login", loginFn},
			{"logoutIntegration", "POST", "/users/logout", logoutFn},
		}

		for _, integration := range integrations {
			apiIntegration, err := apigateway2.Integration(ctx, integration.name, apiGateway, integration.function)
			if err != nil {
				return err
			}

			var route *apigatewayv2.Route
			if integration.name == "loginIntegration" {
				route, err = apigateway2.RouteWithoutAuthorizer(
					ctx,
					fmt.Sprintf(
						"%sRoute",
						integration.name),
					apiGateway,
					integration.httpMethod+" "+integration.path,
					apiIntegration)
			} else {
				route, err = apigateway2.LambdaAuthorizerRoute(
					ctx,
					fmt.Sprintf(
						"%sRoute",
						integration.name),
					apiGateway,
					integration.httpMethod+" "+integration.path,
					apiIntegration,
					lambdaAuthorizer)
			}

			if err != nil {
				return err
			}

			routes = append(routes, route)
		}

		//############################
		// Api Deployment
		apiDeployment, err := apigateway2.Deploy(ctx, "apiDeployment", apiGateway, routes)
		if err != nil {
			errorhandler.HandlingError("Error creating deployment.")
		}

		//############################
		// Api Stage
		stage, err := apigateway2.Stage(ctx, "stage", apiDeployment, apiGateway, logGroup)
		if err != nil {
			errorhandler.HandlingError("Error creating stage.")
		}

		fullApiUrl := pulumi.Sprintf("%s/%s", apiGateway.ApiEndpoint, stage.Name).ApplyT(func(endpoint string) string {
			return endpoint
		}).(pulumi.StringOutput)

		// The URL at which the REST API will be served
		ctx.Export("apiEndpoint", fullApiUrl)

		// Export the User Pool Client ID
		ctx.Export("userPoolClientId", userPoolClient.ID())

		return nil
	})
}
