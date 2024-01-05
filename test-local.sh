#!/usr/bin/env bash
export $(cat .env | xargs)
aws-vault exec cmyk-dev -- go test -v ./handlers/... -run TestStoreAndRetrieveUser
