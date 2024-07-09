import boto3
import header
import json
import logging
import os

logger = logging.getLogger()
logger.setLevel(logging.INFO)

def authorize(event, context):
    token = header.getAuthToken(event)

    try:
        if token:
            if check_token_validity(token):
                logger.info("Token is valid")
                response = generateAllow('user', event['routeArn'])
                return json.loads(response)
            else:
                logger.info("Token is invalid")
                response = generateDeny('user', event['routeArn'])
                return json.loads(response)
    except Exception as e:
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error getting item: {str(e)}')
        }


def check_token_validity(token):
    dynamodb = boto3.client('dynamodb')

    try:
        response = dynamodb.get_item(
            Key={
                'token': {'S': token}
            },
            TableName=os.environ['DYNAMODB_TABLE_NAME'],
        )

        # The token is valid if it exists and 'invalidate' is False
        if 'Item' in response and int(response['Item']['invalidate']['N']) == 0:
            return True
        return False
    except Exception as e:
        logger.error(f"Error fetching item from DynamoDB: {e}")
        return False


def generatePolicy(principalId, effect, resource):
    authResponse = {}
    authResponse['principalId'] = principalId
    if (effect and resource):
        policyDocument = {}
        policyDocument['Version'] = '2012-10-17'
        policyDocument['Statement'] = []
        statementOne = {}
        statementOne['Action'] = 'execute-api:Invoke'
        statementOne['Effect'] = effect
        statementOne['Resource'] = resource
        policyDocument['Statement'] = [statementOne]
        authResponse['policyDocument'] = policyDocument

    authResponse_JSON = json.dumps(authResponse)

    return authResponse_JSON


def generateAllow(principalId, resource):
    return generatePolicy(principalId, 'Allow', resource)


def generateDeny(principalId, resource):
    return generatePolicy(principalId, 'Deny', resource)