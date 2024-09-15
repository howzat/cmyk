#!/usr/bin/env bash
export $(cat .env.local | xargs)
go test ./handlers/... -run TestStoreAndRetrieveUser
