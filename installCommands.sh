#!/usr/bin/env bash

echo "INSTALL COMMANDS"
npm ci
npm run exportEnv
npm run sls --version
make build