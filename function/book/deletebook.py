import boto3
import os
import json

def deleteBook(event, context):
    dynamodb = boto3.client('dynamodb')

    try:
        isbn = event['pathParameters']['isbn']

        # Perform transactional write to delete items from both tables
        dynamodb.transact_write_items(
            TransactItems=[
                {
                    'Delete': {
                        'Key': {
                            'isbn': {
                                'S': isbn,
                            },
                        },
                        'TableName': os.environ['DYNAMODB_TABLE_NAME']
                    }
                },
                {
                    'Delete': {
                        'Key': {
                            'isbn': {
                                'S': isbn,
                            },
                        },
                        'TableName': os.environ['UPDATE_TABLE_NAME']
                    }
                }
            ]
        )

        return {
            "headers": {
                "Content-Type": "application/json"
            },
            'statusCode': 204,
            'body': json.dumps('Book deleted')
        }

    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error deleting item: {str(e)}')
        }
