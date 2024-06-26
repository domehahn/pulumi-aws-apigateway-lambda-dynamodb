import json
import os
import boto3
import body


def createBook(event, context):
    dynamodb = boto3.client('dynamodb')

    try:
        b = body.getBody(event)

        if b:
            author = b.get('author')
            title = b.get('title')
            price = b.get('price')
            isbn = b.get('isbn')
            copiesInStock = b.get('copiesInStock')

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
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 201,
                'message': 'Book created',
                'body': json.dumps(b)
            }

    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error creating item: {str(e)}')
        }
