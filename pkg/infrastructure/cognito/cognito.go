package cognito

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cognito"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func UserPool(ctx *pulumi.Context, name string) (*cognito.UserPool, error) {
	userPool, err := cognito.NewUserPool(ctx, name, &cognito.UserPoolArgs{
		Name: pulumi.String(name),
		AutoVerifiedAttributes: pulumi.StringArray{
			pulumi.String("email"),
		},
		Schemas: cognito.UserPoolSchemaArray{
			&cognito.UserPoolSchemaArgs{
				AttributeDataType: pulumi.String("String"),
				Name:              pulumi.String("email"),
				Mutable:           pulumi.Bool(false),
				Required:          pulumi.Bool(true),
			},
		},
		PasswordPolicy: &cognito.UserPoolPasswordPolicyArgs{
			MinimumLength:                 pulumi.Int(8),
			RequireLowercase:              pulumi.Bool(true),
			RequireNumbers:                pulumi.Bool(true),
			RequireSymbols:                pulumi.Bool(true),
			RequireUppercase:              pulumi.Bool(true),
			TemporaryPasswordValidityDays: pulumi.Int(7),
		},
	})
	return userPool, err
}

func UserPoolClient(ctx *pulumi.Context, name string, userPool *cognito.UserPool) (*cognito.UserPoolClient, error) {
	userPoolClient, err := cognito.NewUserPoolClient(ctx, name, &cognito.UserPoolClientArgs{
		UserPoolId: userPool.ID(),
		ExplicitAuthFlows: pulumi.StringArray{
			pulumi.String("ALLOW_USER_PASSWORD_AUTH"),
			pulumi.String("ALLOW_REFRESH_TOKEN_AUTH"),
		},
		GenerateSecret: pulumi.Bool(false),
	})
	return userPoolClient, err
}
