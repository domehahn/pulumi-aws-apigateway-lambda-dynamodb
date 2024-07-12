import boto3
import os
import json

def getBooks(event, context):
    dynamodb = boto3.client('dynamodb')

    try:
        response = dynamodb.scan(
            ExpressionAttributeNames={
                '#A': 'author',
                '#T': 'title',
                '#P': 'price',
                '#I': 'isbn',
                '#C': 'copiesInStock'
            },
            ProjectionExpression='#A, #T, #P, #I, #C',
            TableName=os.environ['DYNAMODB_TABLE_NAME'],
        )

        if response:
            return {
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 200,
                'message': 'A list of book',
                'body': json.dumps(response["Items"])
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