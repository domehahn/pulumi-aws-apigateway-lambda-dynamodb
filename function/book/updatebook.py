import boto3
import json
import os
from function import body

def updateBook(event, context):
    dynamodb = boto3.client('dynamodb')

    b = body.getBody(event)

    if b:
        try:
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
                        'S': b.get('isbn'),
                    },
                },
                ReturnValues='ALL_NEW',
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
                UpdateExpression='SET #A = :a, #T = :t, #P = :p, #C = :c',
            )

            return {
                'statusCode': 200,
                'body': json.dumps('Data successful updated.')
            }

        except Exception as e:
            return {
                'statusCode': 500,
                'body': json.dumps(f'Error updating item: {str(e)}')
            }

    return {
        'statusCode': 400,
        'body': json.dumps('Data update failed.')
    }