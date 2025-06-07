import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import path = require("path");
import { Platform } from "aws-cdk-lib/aws-ecr-assets";
import * as logs from "aws-cdk-lib/aws-logs";
import * as iam from "aws-cdk-lib/aws-iam";

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    const smPolicy = new iam.PolicyStatement({
      actions: [
        "secretsmanager:GetSecretValue",
        "secretsmanager:DescribeSecret",
      ],
      resources: [
        `arn:aws:secretsmanager:${this.region}:${this.account}:secret:prod/urlshort/redis`, 
      ],
    });
    const func = new lambda.DockerImageFunction(this, "urlShortLambda", {
      functionName: "urlShortLambda",
      code: lambda.DockerImageCode.fromImageAsset(
        path.join(__dirname, "../../"),
        {
          platform: Platform.LINUX_AMD64
        }
      ),
      logRetention: logs.RetentionDays.FOUR_MONTHS,
      timeout: cdk.Duration.seconds(30),
      initialPolicy: [smPolicy]
    });

    // Create API Gateway
    const api = new apigateway.LambdaRestApi(this, "urlShortApi", {
      handler: func,
      proxy: true,
      deployOptions: {
        stageName: "prod",
      },
      defaultCorsPreflightOptions: {
        allowOrigins: apigateway.Cors.ALL_ORIGINS,
        allowMethods: apigateway.Cors.ALL_METHODS,
        allowHeaders: ["Content-Type", "X-Amz-Date", "Authorization", "X-Api-Key", "X-Amz-Security-Token"],
      },
    });

    // Output the API Gateway URL
    new cdk.CfnOutput(this, "ApiUrl", {
      value: api.url,
      description: "URL Shortener API Gateway endpoint",
    });
  }
}
