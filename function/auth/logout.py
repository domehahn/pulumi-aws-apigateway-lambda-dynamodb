import json
import header
import boto3
import os


def logout(event, context):
    client = boto3.client('cognito-idp')

    try:
        token = header.getAuthToken(event)

        if token:

            client.global_sign_out(
                AccessToken=token
            )

            invalidate_token(token)

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


def invalidate_token(token):
    dynamodb = boto3.client('dynamodb')

    dynamodb.put_item(
        Item={
            'token': {'S': token},
            'invalidate': {'N': '1'}
        },
        TableName=os.environ['DYNAMODB_TABLE_NAME']
    )
