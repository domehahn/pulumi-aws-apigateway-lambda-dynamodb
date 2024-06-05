import json
import header
import boto3

def logout(event, context):
    client = boto3.client('cognito-idp')

    try:
        auth = header.getAuthToken(event)

        if auth:
            client.global_sign_out(
                AccessToken=auth.get_token()
            )
            return {
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 200,
                'message': 'Logout successful',
                'body': json.dumps('Logout was successful. User is logged out.')
            }
        return {
            'statusCode': 401,
            'body': json.dumps('Not logged out.')
        }
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error getting item: {str(e)}')
        }