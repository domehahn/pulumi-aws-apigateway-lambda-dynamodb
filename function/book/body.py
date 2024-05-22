import base64
import json

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