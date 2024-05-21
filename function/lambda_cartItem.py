import boto3

def getCartItems(event, context):
    client = boto3.client('dynamodb')
    return client.get_item(
               TableName='book',
               Key={
                   'author': {'S': 'Dan Brown'},
                   'title': {'S': 'Meteor'}
               }
           )