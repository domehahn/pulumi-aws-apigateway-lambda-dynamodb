# AWS Serverless Application with Pulumi
This Project is based on [Pulumi](https://github.com/pulumi/pulumi/tree/master) and written in [Go](https://go.dev/). 
Pulumi offers a Infrastructure as Code SDK for building and deploying infrastructure for several hyper scaler like AWS, 
Azure, etc.

Following services gets created:
1. lambda function
2. dynamodb
5. cloudwatch
6. cognito
7. apigateway

also all necessary roles, policies and attachments.

## Getting Started
The following steps demonstrates how to run the AWS Serverless Application with Pulumi using AWS Serverless Lambda, 
DynamoDb, Apigateway, cloudwatch and Cognito:

1. **Install Pulumi and Golang**:
```bash
$ brew install pulumi/tap/pulumi
$ brew install go
```

2. **Install AWS Cli**
Pulumi makes use of your local AWS Configuration.
```bash
$ brew install awscli
```
After successful installation, you can find the configuration in you home directory.
```bash
$ ${HOME}/.aws
```

2. **Clone GitHub Repository**
```bash
$ mkdir pulumi-playground && cd pulumi-playground
$ git clone https://github.com/domehahn/pulumi-aws-apigateway-lambda-dynamodb.git
```

3. **Configure AWS Cli**
Pulumi needs the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY form you local aws configuration. To configure it you have 
to sign in to AWS.
```bash
$ aws configure sso
```

4. **Configure Pulumi Cli**
```bash
$ pulumi config set aws:profile AWS-Training
$ pulumi config set aws:region eu-central-1
$ pulumi stack select tutorial-playground/pulumi-00/aws-apigateway-lambda-dynamodb
```
Replace `AWS-Training` with the name of your AWS Profile

5. **Run Pulumi**
```bash
$ pulumi up
```

6. **Stop Pulumi**
```bash
$ pulumi destroy
```

### Alternative
For configure AWS and Pulumi you can use the Taskfile within the project. To make use of Taskfile you have to install it.
```bash
$ brew install go-task/tap/go-task
$ brew install go-task
```

1. **Configure AWS Cli**
Inside the project directory run
```bash
$ task aws-configure-sso 
```
or
```bash
$ task aws-configure
```

2. **Configure Pulumi Cli**
```bash
$ task pulumi-config
```

3. **Run Pulumi**
```bash
$ task pulumi-up
```
This runs with the parameter `--yes` which is a auto approval for provisioning the aws services.

4. **Stop Pulumi**
```bash
$ task pulumi-destroy
```
This runs with the parameter `--yes` which is a auto approval for provisioning the aws services.