import json
import boto3
import header
import body
import os
import logging

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def login(event, context):
    client = boto3.client('cognito-idp')
    dynamodb = boto3.client('dynamodb')

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

            if 'AuthenticationResult' in response:
                auth_result = response['AuthenticationResult']
                token = auth_result['AccessToken']

                # Add token to DynamoDB
                dynamodb.put_item(
                    Item={
                        'token': {'S': token},
                        'invalidate': {'N': '0'}
                    },
                    TableName=os.environ['DYNAMODB_TABLE_NAME']
                )

                return {
                    "headers": {
                        "Content-Type": "application/json"
                    },
                    'statusCode': 200,
                    'message': 'Authentication successful',
                    'body': json.dumps(auth_result)
                }

        return {
            'statusCode': 401,
            'body': json.dumps('Not authorized.')
        }
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error during login: {str(e)}')
        }
