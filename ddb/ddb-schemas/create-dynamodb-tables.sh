#!/bin/sh +x

if [ ! $# -eq 1 ]; then
  echo "Usage: sh create-dynamodb-table.sh <table schema folder for all services in container>"
  exit 1
fi

SCHEMA_FOLDER=$1
echo "${SCHEMA_FOLDER}"/*.json
for SCHEMA in "${SCHEMA_FOLDER}"/*.json; do

    jsonfile=$(basename "${SCHEMA##*/}")
    tablename=${jsonfile%.json}
    echo "checking if table ${tablename} exists (using the filename/tablename convention)"
    table_exists=$(aws dynamodb list-tables --endpoint-url "${DYNAMO_ENDPOINT}" --region us-east-1 | jq "(.TableNames | length > 0) and (.TableNames | contains([\"${tablename}\"]))")
    case $table_exists in
      false)
          echo "Creating table from schema ${SCHEMA}"
          aws dynamodb create-table --cli-input-json "file://${SCHEMA}" --endpoint-url "${DYNAMO_ENDPOINT}" --region ${AWS_REGION-eu-west-2} > /dev/null || exit $?
      ;;
      true)
        echo "skip: table from schema file exists ${SCHEMA}"
    esac
done