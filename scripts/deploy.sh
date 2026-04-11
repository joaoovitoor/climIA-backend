#!/bin/bash
set -euo pipefail

FUNCTION_NAME="climia-backend"
REGION="${AWS_REGION:-us-east-1}"
S3_BUCKET="climia-deployments-274616633915"
S3_KEY="${FUNCTION_NAME}/deploy.zip"

echo "=> Building..."
npm run build

echo "=> Packaging..."
rm -rf .deploy
mkdir -p .deploy
cp -r dist .deploy/
cp package.json .deploy/
cd .deploy
npm install --production --ignore-scripts --no-audit --no-fund 2>/dev/null
rm -f package.json package-lock.json
zip -rq deploy.zip .
cd ..

SIZE=$(du -h .deploy/deploy.zip | cut -f1)
echo "=> Package size: ${SIZE}"

echo "=> Uploading to S3..."
aws s3 cp .deploy/deploy.zip "s3://${S3_BUCKET}/${S3_KEY}" --region "$REGION" --quiet

echo "=> Deploying ${FUNCTION_NAME}..."
aws lambda update-function-code \
  --function-name "$FUNCTION_NAME" \
  --s3-bucket "$S3_BUCKET" \
  --s3-key "$S3_KEY" \
  --region "$REGION" \
  --no-cli-pager > /dev/null

echo "=> Waiting for update..."
aws lambda wait function-updated \
  --function-name "$FUNCTION_NAME" \
  --region "$REGION"

echo "=> Done! ${FUNCTION_NAME} deployed."
rm -rf .deploy
