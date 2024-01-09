#!/usr/bin/env bash

export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export AWS_SESSION_TOKEN

echo "INSTALL STEPS ${DEV_ACCOUNT_ID} - ${STAGE}"
CREDS=$(aws sts assume-role --role-arn arn:aws:iam::"${DEV_ACCOUNT_ID}":role/CMYKDeployerRole --role-session-name=CMYKDeployerRole)
AWS_ACCESS_KEY_ID=$(echo $CREDS | jq -r '.Credentials.AccessKeyId')
AWS_SECRET_ACCESS_KEY=$(echo $CREDS | jq -r '.Credentials.SecretAccessKey')
AWS_SESSION_TOKEN=$(echo $CREDS | jq -r '.Credentials.SessionToken')
npm run exportEnv
npm run sls deploy --stage=$STAGE