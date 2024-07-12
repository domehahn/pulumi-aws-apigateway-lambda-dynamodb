#!/bin/bash

# Clear AWS CLI cache
rm -rf ~/.aws/cli/cache/

# Unset environment variables
unset AWS_ACCESS_KEY_ID
unset AWS_SECRET_ACCESS_KEY
unset AWS_SESSION_TOKEN

# Clear shared credentials file
> ~/.aws/credentials

# Refresh Pulumi state
pulumi refresh