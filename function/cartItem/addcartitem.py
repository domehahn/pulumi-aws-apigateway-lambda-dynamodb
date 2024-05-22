import boto3
import body
import json
import os

def addCartItem(event, context):
    dynamodb = boto3.client('dynamodb')

    try:
        b = body.getBody(event)

        if b:
            isbn = b.get('isbn')
            quantity = b.get('quantity')

            dynamodb.put_item(
                Item={
                    'isbn': {
                        'S': isbn

                    },
                    'quantity': {
                        'N': quantity

                    }
                },
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
            )

            return {
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 201,
                'message': 'Item added to cart',
                'body': json.dumps(b)
            }

    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error creating item: {str(e)}')
        }
