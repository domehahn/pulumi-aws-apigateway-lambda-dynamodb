import boto3
import os
import json

def getBook(event, context):
    dynamodb = boto3.client('dynamodb')

    isbn = event['pathParameters']['isbn']

    response = dynamodb.get_item(
        Key={
            'isbn': {
                'S': isbn,
            },
        },
        TableName=os.environ['DYNAMODB_TABLE_NAME'],
    )

    return {
        "headers": {
            "Content-Type": "application/json"
        },
        'statusCode': 200,
        'message': 'A single book',
        'body': json.dumps(response['Item'])
    }