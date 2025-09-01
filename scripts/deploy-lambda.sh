#!/bin/bash

# AWS Lambda deployment script

set -e

FUNCTION_NAME=${1:-users-api}
REGION=${2:-us-east-1}

echo "Deploying $FUNCTION_NAME to $REGION..."

# Clean up previous builds
echo "Cleaning up previous builds..."
rm -f bootstrap lambda.zip

# Build for Linux
echo "Building for Linux/AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap cmd/lambda/main.go

# Create deployment package
echo "Creating deployment package..."
zip lambda.zip bootstrap

# Check if function exists
if aws lambda get-function --function-name "$FUNCTION_NAME" --region "$REGION" >/dev/null 2>&1; then
    echo "Updating existing Lambda function..."
    aws lambda update-function-code \
        --function-name "$FUNCTION_NAME" \
        --zip-file fileb://lambda.zip \
        --region "$REGION"
else
    echo "Function $FUNCTION_NAME does not exist. Please create it first using AWS Console or CLI."
    echo "Example create command:"
    echo "aws lambda create-function \\"
    echo "  --function-name $FUNCTION_NAME \\"
    echo "  --runtime provided.al2 \\"
    echo "  --role arn:aws:iam::YOUR_ACCOUNT:role/lambda-execution-role \\"
    echo "  --handler bootstrap \\"
    echo "  --zip-file fileb://lambda.zip \\"
    echo "  --region $REGION"
    exit 1
fi

echo "Deployment completed successfully!"
echo "Function ARN: $(aws lambda get-function --function-name "$FUNCTION_NAME" --region "$REGION" --query 'Configuration.FunctionArn' --output text)"

# Clean up
rm -f bootstrap lambda.zip