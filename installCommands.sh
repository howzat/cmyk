#!/usr/bin/env bash

echo "INSTALL COMMANDS"
npm ci
npm run sls --version
make build