import json
import header
import boto3
import body

def login(event, context):
    client = boto3.client('cognito-idp')

    try:
        b = body.getBody(event)
        auth = header.getBasicAuthValues(event)

        if auth and b:
            response = client.initiate_auth(
                AuthFlow='USER_PASSWORD_AUTH',
                AuthParameters={
                    'USERNAME': auth.get_username(),
                    'PASSWORD': auth.get_password()
                },
                ClientId=b.get('clientId')
            )

            if response:
                return {
                    "headers": {
                        "Content-Type": "application/json"
                    },
                    'statusCode': 200,
                    'message': 'Authentication successful',
                    'body': json.dumps(response['AuthenticationResult'])
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