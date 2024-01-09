#!/usr/bin/env bash
set -x
echo "INSTALL STEPS ${{DEV_ACCOUNT_ID}}"
CREDS=`aws sts assume-role --role-arn arn:aws:iam::${{DEV_ACCOUNT_ID}}:role/CMYKDeployerRole --role-session-name=CMYKDeployerRole`
export AWS_ACCESS_KEY_ID=`echo $CREDS | jq -r '.Credentials.AccessKeyId'`
export AWS_SECRET_ACCESS_KEY=`echo $CREDS | jq -r '.Credentials.SecretAccessKey'`
export AWS_SESSION_TOKEN=`echo $CREDS | jq -r '.Credentials.SessionToken'`
npm run sls deploy --stage=$STAGE