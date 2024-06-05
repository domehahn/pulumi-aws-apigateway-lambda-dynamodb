package cognito

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cognito"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SensitiveUser struct {
	Username string
	Password string
	Email    string
}

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

func User(ctx *pulumi.Context, name string, sensitiveUser SensitiveUser, userPool *cognito.UserPool) (*cognito.User, error) {
	user, err := cognito.NewUser(ctx, name, &cognito.UserArgs{
		UserPoolId: userPool.ID(),
		Username:   pulumi.String(sensitiveUser.Username),
		Password:   pulumi.String(sensitiveUser.Password), // Use a strong password
		Attributes: pulumi.StringMap{
			"email": pulumi.String(sensitiveUser.Email),
		},
	})
	return user, err
}

func UserPoolClient(ctx *pulumi.Context, name string, userPool *cognito.UserPool) (*cognito.UserPoolClient, error) {
	userPoolClient, err := cognito.NewUserPoolClient(ctx, name, &cognito.UserPoolClientArgs{
		UserPoolId: userPool.ID(),
		ExplicitAuthFlows: pulumi.StringArray{
			pulumi.String("ALLOW_USER_PASSWORD_AUTH"),
			pulumi.String("ALLOW_REFRESH_TOKEN_AUTH"),
		},
		GenerateSecret: pulumi.Bool(false),
	}, pulumi.DependsOn([]pulumi.Resource{userPool}))
	return userPoolClient, err
}
