import boto3
import json
import os
import body


def updateBook(event, context):
    dynamodb = boto3.client('dynamodb')

    try:
        b = body.getBody(event)

        isbn = event['pathParameters']['isbn']

        if b:
            dynamodb.update_item(
                ExpressionAttributeNames={
                    '#A': 'author',
                    '#T': 'title',
                    '#P': 'price',
                    '#C': 'copiesInStock',
                },
                ExpressionAttributeValues={
                    ':a': {
                        'S': b.get('author'),
                    },
                    ':t': {
                        'S': b.get('title'),
                    },
                    ':p': {
                        'N': str(b.get('price')),
                    },
                    ':c': {
                        'N': str(b.get('copiesInStock')),
                    },
                },
                Key={
                    'isbn': {
                        'S': isbn,
                    },
                },
                ReturnValues='ALL_NEW',
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
                UpdateExpression='SET #A = :a, #T = :t, #P = :p, #C = :c',
            )

            return {
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 200,
                'message': 'Book updated',
                'body': json.dumps(b)
            }

        return {
            'statusCode': 400,
            'body': json.dumps('Data update failed.')
        }

    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error updating item: {str(e)}')
        }
