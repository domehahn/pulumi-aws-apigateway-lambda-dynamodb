import json
import os
import boto3
from function import body


def createBook(event, context):
    dynamodb = boto3.client('dynamodb')

    b = body.getBody(event)

    if b:
        author = b.get('author')
        title = b.get('title')
        price = b.get('price')
        isbn = b.get('isbn')
        copiesInStock = b.get('copiesInStock')

        try:
            dynamodb.put_item(
                Item={
                    'author': {
                        'S': author

                    },
                    'title': {
                        'S': title

                    },
                    'price': {
                        'N': price
                    },
                    'isbn': {
                        'S': isbn

                    },
                    'copiesInStock': {
                        'N': copiesInStock

                    }
                },
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
            )

            return {
                'statusCode': 200,
                'body': json.dumps('Data successful created.')
            }

        except Exception as e:
            return {
                'statusCode': 500,
                'body': json.dumps(f'Error creating item: {str(e)}')
            }

    return {
        'statusCode': 400,
        'body': json.dumps('Data creation failed.')
    }