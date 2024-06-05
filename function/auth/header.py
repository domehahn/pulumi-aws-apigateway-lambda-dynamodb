import base64
import json
import logging


class BasicAuth:
    def __init__(self, username, password):
        self.__username = username  # Private attribute
        self.__password = password  # Private attribute

    # Getter method for username
    def get_username(self):
        return self.__username

    # Getter method for password
    def get_password(self):
        return self.__password


class Token:
    def __init__(self, token):
        self.__token = token

    # Getter method for token
    def get_token(self):
        return self.__token


logger = logging.getLogger("boto3")
logger.setLevel(logging.INFO)


def getBasicAuthValues(event):
    headers = event.get('headers', {})

    auth_header = headers.get('authorization')
    if auth_header is None or not auth_header.startswith('Basic '):
        return {
            'statusCode': 401,
            'body': json.dumps('Unauthorized')
        }

    if auth_header:
        encoded_credentials = auth_header.split(' ')[1]
        decoded_credentials = base64.b64decode(encoded_credentials).decode('utf-8')
        username, password = decoded_credentials.split(':', 1)
        auth = BasicAuth(username, password)
        return auth


def getAuthToken(event):
    headers = event.get('headers', {})

    auth_header = headers.get('authorization')
    if auth_header is None or not auth_header.startswith('Bearer '):
        return {
            'statusCode': 401,
            'body': json.dumps('Unauthorized')
        }

    if auth_header:
        # Extract the token part after 'Bearer '
        bearer_token = auth_header[len('Bearer '):]
        auth = Token(bearer_token)
        return auth
