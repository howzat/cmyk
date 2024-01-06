#!/usr/bin/env bash
set -x
echo "INSTALL COMMANDS"
npm ci
npm run exportEnv
npm run sls --version
make build