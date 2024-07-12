import boto3
import os
import json

def deleteCartItem(event, context):
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

        # Check if the item exists and retrieve the quantity
        if 'Item' not in response or 'quantity' not in response['Item']:
            return {
                'statusCode': 404,
                'body': json.dumps(f'Item with ISBN {isbn} not found in the cart.')
            }

        quantity_to_add = int(response['Item']['quantity']['N'])

        # Perform transactional write to delete item from cart and update product stock
        dynamodb.transact_write_items(
            TransactItems=[
                {
                    'Delete': {
                        'Key': {
                            'isbn': {'S': isbn},
                        },
                        'TableName': os.environ['DYNAMODB_TABLE_NAME']  # Cart table
                    }
                },
                {
                    'Update': {
                        'TableName': os.environ['UPDATE_TABLE_NAME'],  # Products table
                        'Key': {
                            'isbn': {'S': isbn},
                        },
                        'UpdateExpression': 'SET copiesInStock = copiesInStock + :q',
                        'ExpressionAttributeValues': {
                            ':q': {'N': str(quantity_to_add)},
                        }
                    }
                }
            ]
        )

        return {
            "headers": {
                "Content-Type": "application/json"
            },
            'statusCode': 204,
            'body': json.dumps('Item removed from cart')
        }

    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error deleting item: {str(e)}')
        }