import boto3
import os



def getBooks(event, context):
    dynamodb = boto3.resource('dynamodb')
    table = dynamodb.Table(os.environ['DYNAMODB_TABLE_NAME'])
    response = table.scan()
    return response['Items']

def createBook(event, context):
    dynamodb = boto3.client('dynamodb')
    dynamodb.put_item(
        Item={
            'id': {
              'S': '1'
            },
            'author': {
               'S': 'Dan Brown'

            },
            'title': {
                'S': 'Meteor'

            },
            'price': {
                'N': '14.00'
            },
            'isbn': {
                'S': '978-3404175048'

            },
            'copiesInStock': {
                'N': '800'

            }
        },
        TableName=os.environ['DYNAMODB_TABLE_NAME'],
    )