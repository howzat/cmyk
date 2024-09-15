#!/usr/bin/env bash
export $(cat .env | xargs)
aws-vault exec cmyk-dev -- go test ./handlers/... -run TestStoreAndRetrieveUser
