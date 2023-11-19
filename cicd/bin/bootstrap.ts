#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {GithubSourceConfig, PipelineConfig, WebsitePipeline} from "../lib/pipeline";

const app = new cdk.App();
const appName = "CMYK"
const sourceOwner = app.node.tryGetContext('sourceOwner');
const sourceRepositoryName = app.node.tryGetContext('sourceRepository');
const sourceConfig = new GithubSourceConfig("main", sourceOwner, sourceRepositoryName);
const pipelineConfig = new PipelineConfig(appName, sourceConfig);

let env = {
    account: process.env.CDK_DEPLOY_ACCOUNT || process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEPLOY_REGION || process.env.CDK_DEFAULT_REGION
};


new WebsitePipeline(
  app,
  'WebsitePipelineStack',
  pipelineConfig,
  {env}
)


app.synth();