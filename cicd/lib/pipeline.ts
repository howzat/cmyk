import {SecretValue, Stack, StackProps} from 'aws-cdk-lib';
import {Construct} from 'constructs';
import {GitHubTrigger} from "aws-cdk-lib/aws-codepipeline-actions";
import {CodePipeline, CodePipelineSource, ShellStep} from "aws-cdk-lib/pipelines";
import {StaticWebsiteHostingDeployStage} from "./deploy-stage";
import {BuildSpec} from "aws-cdk-lib/aws-codebuild";

export class GithubSourceConfig {
    constructor(branchName: string, owner: string, repositoryName: string) {
        this.branchName = branchName
        this.owner = owner
        this.repositoryName = repositoryName;
    }

    public readonly branchName: string
    public readonly owner: string
    public readonly repositoryName: string;

    public repositoryString(): string {
        return `${this.owner}/${this.repositoryName}`
    }

    public gitHubTokenPath(): string {
        return `/github/${this.repositoryString()}`
    }
}

export class PipelineConfig {
    constructor(appName: string, githubSourceConfig: GithubSourceConfig) {
        this.appName = appName;
        this.githubSourceConfig = githubSourceConfig
    }

    public readonly githubSourceConfig: GithubSourceConfig;
    public readonly appName: string;
}

export class WebsitePipeline extends Stack {
    private readonly config: PipelineConfig;
    private readonly authentication: SecretValue;

    constructor(scope: Construct, id: string, config: PipelineConfig, props?: StackProps) {
        super(scope, id, props);

        this.config = config;

        //CodePipeline object
        let githubConfig = config.githubSourceConfig;
        let repoString = githubConfig.repositoryString();
        console.log("repoString", repoString)
        console.log("branchName", githubConfig.branchName)
        console.log("githubTokenPath", githubConfig.gitHubTokenPath())

        this.authentication = SecretValue.secretsManager(githubConfig.gitHubTokenPath());
        console.log("authentication", this.authentication)

        let githubSource = CodePipelineSource.gitHub(repoString, githubConfig.branchName, {
            trigger: GitHubTrigger.WEBHOOK,
            authentication: this.authentication,
        });

        const codePipeline = new CodePipeline(this, `${config.appName}BuildPipeline`, {
            pipelineName: 'CMYKCICDPipeline',
            crossAccountKeys: false,
            codeBuildDefaults: {
                partialBuildSpec: BuildSpec.fromObject({
                    version: "0.2",
                    phases: {
                        install: {
                            "runtime-versions": {
                                nodejs: 14
                            },
                            commands: [
                                "n 18.2.0",
                            ]
                        }
                    }
                })
            },
            synth: new ShellStep('Synth', {
                primaryOutputDirectory: 'cicd/cdk.out',
                input: githubSource,
                installCommands: [
                    "./build.sh"
                ],
                commands: [
                    "cd ./cicd",
                    "npx cdk synth",
                ],
            }),
        });

        // codePipeline.addStage(new ApprovalStage(this, 'ApprovalStage'), {});

        let staticWebsiteHostingDeployStage = new StaticWebsiteHostingDeployStage(this, 'DeployStage', {env: props?.env});
        codePipeline.addStage(staticWebsiteHostingDeployStage, {})
    }
}