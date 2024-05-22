import boto3
import os
import json
from function import body

def deleteBook(event, context):
    dynamodb = boto3.client('dynamodb')

    b = body.getBody(event)

    if b:
        try:
            dynamodb.delete_item(
                Key={
                    'isbn': {
                        'S': b.get('isbn'),
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