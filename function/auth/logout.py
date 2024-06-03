import json

import boto3

def logout(event, context):
    client = boto3.client('cognito-idp')

    try:
        b = body.getBody(event)

        if b:
            client.global_sign_out(
                AccessToken=b.get('accessToken'),
                ClientId=b.get('clientId')
            )
            return {
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 200,
                'message': 'Logout successful',
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