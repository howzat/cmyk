# CMYK CICD IaC

The project directory contains both the code that creates the code building pipeline and the code the pipeline is
going to build and deploy. The CICD `bootstrap.ts` CDK project creates the CodeBuild Pipeline that builds and 
deploys a Gatsby static site to an S3 bucket with a Cloudfront distribution.

To provision the pipeline make sure the target environment has been bootstrapped
```shell
cdk bootstrap aws://ACCOUNT-NUMBER-1/REGION-1
```

Then execute deploying the WebsitePipelineStack once. Afterward the pipeline will keep itself up-to-date.
[See the Pipelines CDK construct for more information](https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.pipelines-readme.html)

>Important: be sure to git commit and git push before deploying the Pipeline stack using cdk deploy!
>The reason is that the pipeline will start deploying and self-mutating right away based on the sources in the repository, so the sources it finds in there should be the ones you want it to find.

# Secrets used in building the Pipeline
The pipeline relies on a Github token stored (by hand) in AWS `SecretsManager` to register hooks for triggering the build.
Running the script should produce somthing like the following output.

This should only need to be done when a change to the pipeline is required.

```shell
$ cd cicd
$ git commit -a
$ git push
$ cdk deploy WebsitePipelineStack
$ aws-vault exec tsdev -- cdk deploy -c sourceOwner=howzat -c sourceRepository=cmyk
Enter token for arn:aws:iam::XXXXX:mfa/XXX: 926768
SecretValue {
  creationStack: [ 'stack traces disabled' ],
  value: CfnDynamicReference {
    creationStack: [ 'stack traces disabled' ],
    value: '{{resolve:secretsmanager:/github/howzat/cmyk:SecretString:::}}'
  },
  rawValue: CfnDynamicReference {
    creationStack: [ 'stack traces disabled' ],
    value: '{{resolve:secretsmanager:/github/howzat/cmyk:SecretString:::}}'
  }
}
${Token[TOKEN.198]}

✨  Synthesis time: 7.23s

WebsitePipelineStack: deploying...
[0%] start: Publishing c21926fa4cbb9eb34d01c33edecc868ba74f072ad695a87ba2012aa3fee3acff:678975692412-eu-west-1
[100%] success: Published c21926fa4cbb9eb34d01c33edecc868ba74f072ad695a87ba2012aa3fee3acff:678975692412-eu-west-1
WebsitePipelineStack: creating CloudFormation changeset...

 ✅  WebsitePipelineStack

✨  Deployment time: 33.95s
```

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful CDK commands

* `npm run build`   compile typescript to js
* `npm run watch`   watch for changes and compile
* `npm run test`    perform the jest unit tests
* `cdk deploy`      deploy this stack to your default AWS account/region
* `cdk diff`        compare deployed stack with current state
* `cdk synth`       emits the synthesized CloudFormation template
