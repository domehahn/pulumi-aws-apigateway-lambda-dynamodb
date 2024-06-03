import json
import body
import boto3

def login(event, context):
    client = boto3.client('cognito-idp')

    try:
        b = body.getBody(event)

        if b:
            response = client.initiate_auth(
                AuthFlow='USER_PASSWORD_AUTH',
                AuthParameters={
                    'USERNAME': b.get('username'),
                    'PASSWORD': b.get('password')
                },
                ClientId=b.get('clientId'),
                UserPoolId=b.get('userPoolId')
            )

            if response:
                return {
                    "headers": {
                        "Content-Type": "application/json"
                    },
                    'statusCode': 200,
                    'message': 'Authentication successful',
                    'body': json.dump(response)
                }

            return {
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 200,
                'message': 'No book found'
            }

        return {
            'statusCode': 401,
            'body': json.dumps('Not authorized.')
        }
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error getting item: {str(e)}')
        }