#!/usr/bin/env bash
set -e
set -u
set -o pipefail

export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export AWS_SESSION_TOKEN
export SLS_STAGE

echo "Assuming role CMYKDeployerRole in account [${DEV_ACCOUNT_ID}] - Stage: [${SLS_STAGE}] - Region: [${AWS_REGION}]"

if [ -z "${DEV_ACCOUNT_ID}" ]; then
  >&2 echo Error: DEV_ACCOUNT_ID must be set. Assuming cross account role has failed!
  exit 1
else
  CREDS=$(aws sts assume-role --role-arn arn:aws:iam::"${DEV_ACCOUNT_ID}":role/CMYKDeployerRole --role-session-name=CMYKDeployerRole)
  AWS_ACCESS_KEY_ID=$(echo $CREDS | jq -r '.Credentials.AccessKeyId')
  AWS_SECRET_ACCESS_KEY=$(echo $CREDS | jq -r '.Credentials.SecretAccessKey')
  AWS_SESSION_TOKEN=$(echo $CREDS | jq -r '.Credentials.SessionToken')
  npm run exportEnv
  npm run sls deploy --stage="${SLS_STAGE}" --verbose --region="${AWS_REGION}"
fi

