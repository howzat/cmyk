import {AllowedMethods, Distribution, SecurityPolicyProtocol, ViewerProtocolPolicy} from "aws-cdk-lib/aws-cloudfront";

import {aws_s3 as S3, CfnOutput, Duration, RemovalPolicy, Stack, StackProps} from 'aws-cdk-lib';
import {BlockPublicAccess, BucketAccessControl} from "aws-cdk-lib/aws-s3";
import {S3Origin} from "aws-cdk-lib/aws-cloudfront-origins";
import {WebsiteCertificatesStack} from "./website-certificates";
import {ARecord, RecordTarget} from "aws-cdk-lib/aws-route53";
import {CloudFrontTarget} from "aws-cdk-lib/aws-route53-targets";
import {Construct} from "constructs";
import {BucketDeployment, Source} from "aws-cdk-lib/aws-s3-deployment";
import * as path from "path";
import * as fs from 'fs';

export interface StaticWebsiteStackProps extends StackProps {
    apexDomain: string;
    webSubdomain: string;
    webCerts: WebsiteCertificatesStack
}

export class StaticWebsiteStack extends Stack {

    constructor(parent: Construct, name: string, props: StaticWebsiteStackProps) {
        super(parent, name, props);

        let wwwDomain = `${props.webSubdomain}.${props.apexDomain}`;
        console.log("wwwDomain is " + wwwDomain)
        // The behaviour of the distribution is dynamic see: https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_cloudfront-readme.html
        // This bucket is handled as a bucket origin and CloudFront's redirect and error handling will be used.
        // The Origin will create an origin access identity and grant it access to the underlying bucket.
        // This can be used in conjunction with a bucket that is not public to require that your users access your content using CloudFront URLs and not S3 URLs directly.
        const website = new S3.Bucket(this, 'WebsiteBucket', {
            bucketName: `site-${props.apexDomain}`,
            autoDeleteObjects: true,
            removalPolicy: RemovalPolicy.DESTROY,
            publicReadAccess: false,
            accessControl: BucketAccessControl.PRIVATE,
            blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
        })

        console.log("content bucket created", website.bucketArn)

        const distribution = new Distribution(this, 'CloudFrontDistribution', {
            domainNames: [props.apexDomain, wwwDomain],
            defaultRootObject: 'index.html',
            certificate: props.webCerts.siteCertificate,
            minimumProtocolVersion: SecurityPolicyProtocol.TLS_V1_2_2021,
            defaultBehavior: {
                compress: true,
                origin: new S3Origin(website),
                allowedMethods: AllowedMethods.ALLOW_ALL,
                viewerProtocolPolicy: ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
            },
            errorResponses: [
                {
                    httpStatus: 403,
                    responseHttpStatus: 403,
                    responsePagePath: '/404.html',
                    ttl: Duration.minutes(30),
                }
            ],
        });

        new ARecord(this, 'SiteAliasRecord', {
            recordName: wwwDomain,
            target: RecordTarget.fromAlias(new CloudFrontTarget(distribution)),
            zone: props.webCerts.hostedZone
        });

        let contentsLocation = path.join("./", "..", "webapp", "cmyk", "public");
        let resolve = path.resolve(contentsLocation);
        console.log("resolved location of assets contains", getAllFilesFromFolder(resolve))

        let bucketDeployment = new BucketDeployment(this, 'DeployWebsite', {
            sources: [Source.asset(contentsLocation)],
            destinationBucket: website,
            distributionPaths: ['/*'],
            distribution: distribution,
        });

        bucketDeployment.node.addDependency(website, distribution)

        new CfnOutput(this, "WebsiteContentBucketDomainName", {value: website.bucketName,});
        new CfnOutput(this, "WebsiteContentBucketName", {value: website.bucketDomainName,});
        new CfnOutput(this, 'WebsiteContentBucket', {value: website.bucketArn});
        new CfnOutput(this, 'DistributionDomainName', {value: distribution.domainName});
        new CfnOutput(this, 'DistributionId', {value: distribution.distributionId});
    }
}

const getAllFilesFromFolder = (dir: string) => {

    let results: string[] = [];
    fs.readdirSync(dir).forEach((file: string) => {

        file = dir + '/' + file;
        let stat = fs.statSync(file);

        if (stat && stat.isDirectory()) {
            results = results.concat(getAllFilesFromFolder(file))
        } else results.push(file);

    });

    return results;
};
