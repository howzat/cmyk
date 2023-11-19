#!/usr/bin/env bash

pushd .  >> /dev/null
cd ./cicd || exit
npm ci
npm install -g aws-cdk
popd || exit

pushd .  >> /dev/null
cd ./cmyk.api || exit
popd || exit

pushd .  >> /dev/null
cd ./cmyk.webapp || exit
popd || exit