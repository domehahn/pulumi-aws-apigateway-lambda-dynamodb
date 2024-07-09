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
            quantity = int(b.get('quantity'))

            # Prepare the transaction items
            transact_items = [
                {
                    'Update': {
                        'TableName': os.environ['UPDATE_TABLE_NAME'],
                        'Key': {
                            'isbn': {
                                'S': isbn,
                            },
                        },
                        'UpdateExpression': 'SET #C = #C - :q',
                        'ExpressionAttributeNames': {
                            '#C': 'copiesInStock',
                        },
                        'ExpressionAttributeValues': {
                            ':q': {
                                'N': str(quantity),
                            },
                        },
                        'ConditionExpression': '#C >= :q',
                    }
                }
            ]

            # Check if item already exists in the cart
            response = dynamodb.get_item(
                Key={
                    'isbn': {
                        'S': isbn,
                    },
                },
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
            )

            if 'Item' not in response or 'quantity' not in response['Item']:
                # Item does not exist in the cart, initialize quantity
                new_quantity = quantity
            else:
                # Item exists in the cart, get current quantity
                quantity_in_cart = int(response['Item']['quantity']['N'])
                new_quantity = quantity_in_cart + quantity

            transact_items.append({
                'Update': {
                    'TableName': os.environ['DYNAMODB_TABLE_NAME'],
                    'Key': {
                        'isbn': {
                            'S': isbn,
                        },
                    },
                    'UpdateExpression': 'SET #Q = :q',
                    'ExpressionAttributeNames': {
                        '#Q': 'quantity',
                    },
                    'ExpressionAttributeValues': {
                        ':q': {
                            'N': str(new_quantity),
                        },
                    }
                }
            })

            # Perform the transaction
            dynamodb.transact_write_items(TransactItems=transact_items)

            return {
                "headers": {
                    "Content-Type": "application/json"
                },
                'statusCode': 201,
                'message': 'Item added to cart',
                'body': json.dumps(b)
            }

    except dynamodb.exceptions.TransactionCanceledException as e:
        return {
            'statusCode': 400,
            'body': json.dumps('Item(s) could not be added to cart. Copies in stock are less than quantity.'),
        }
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error creating item: {str(e)}')
        }
