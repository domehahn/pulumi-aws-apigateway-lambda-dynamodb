import boto3
import os
import json


def getBooks(event, context):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(os.environ['DYNAMODB_TABLE_NAME'])
    response = table.scan()
    return {
        "headers": {
            "Content-Type": "application/json"
        },
        'statusCode': 200,
        'message': 'A list of book',
        'body': response['Items']
    }
