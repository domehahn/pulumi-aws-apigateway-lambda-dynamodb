import boto3
import os
import json
import uuid
import base64


def getBooks(event, context):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(os.environ['DYNAMODB_TABLE_NAME'])
    response = table.scan()
    return response['Items']


def createBook(event, context):
    dynamodb = boto3.client('dynamodb')

    body = getBody(event)

    if body:
        author = body.get('author')
        title = body.get('title')
        price = body.get('price')
        isbn = body.get('isbn')
        copiesInStock = body.get('copiesInStock')

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

def updateBook(event, context):
    dynamodb = boto3.client('dynamodb')

    body = getBody(event)

    if body:
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
                        'S': body.get('author'),
                    },
                    ':t': {
                        'S': body.get('title'),
                    },
                    ':p': {
                        'N': str(body.get('price')),
                    },
                    ':c': {
                        'N': str(body.get('copiesInStock')),
                    },
                },
                Key={
                    'isbn': {
                        'S': body.get('isbn'),
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

def deleteBook(event, context):
    dynamodb = boto3.client('dynamodb')

    body = getBody(event)

    if body:
        try:
            dynamodb.delete_item(
                Key={
                    'isbn': {
                        'S': body.get('isbn'),
                    },
                },
                TableName=os.environ['DYNAMODB_TABLE_NAME'],
            )
            return {
                'statusCode': 200,
                'body': json.dumps('Data successful deleted.')
            }

        except Exception as e:
            return {
                'statusCode': 500,
                'body': json.dumps(f'Error deleting item: {str(e)}')
            }

    return {
        'statusCode': 400,
        'body': json.dumps('Data deletion failed.')
    }

def getBody(event):
    is_base64_encoded = event.get('isBase64Encoded', False)
    body = event.get('body')

    if is_base64_encoded and body:
        body = base64.b64decode(body)

    # If the body is JSON, parse it
    if isinstance(body, bytes):
        body = body.decode('utf-8')

    try:
        body = json.loads(body)
    except json.JSONDecodeError:
        return {
            'statusCode': 400,
            'body': json.dumps('Invalid JSON')
        }
    return body