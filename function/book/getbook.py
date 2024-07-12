import boto3
import os
import json

def getBook(event, context):
    dynamodb = boto3.client('dynamodb')

    try:
        isbn = event['pathParameters']['isbn']

        response = dynamodb.get_item(
            Key={
                'isbn': {
                    'S': isbn,
                },
            },
            TableName=os.environ['DYNAMODB_TABLE_NAME'],
        )

        if response:
            return {
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 200,
                'message': 'A single book',
                'body': json.dumps(response["Item"])
            }

        return {
            "headers": {
                "Content-Type": "application/json"
            },
            'statusCode': 200,
            'message': 'No book found'
        }

    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error getting item: {str(e)}')
        }