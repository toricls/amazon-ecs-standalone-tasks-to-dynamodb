# amazon-ecs-standalone-tasks-to-dynamodb

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[license]: https://github.com/toricls/amazon-ecs-standalone-tasks-to-dynamodb/blob/master/LICENSE
A sample AWS SAM application to store `RUNNING` ECS standalone tasks metadata in a DynamoDB table using EventBridge Events and a Lambda function.

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Build and Install

### Installing dependencies & building the target

```shell
$ sam build
```

### Packaging and deployment

```bash
$ sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
* **AWS Region**: The AWS region you want to deploy your app to.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
* **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modifies IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

## Contribution

1. Fork ([https://github.com/toricls/amazon-ecs-standalone-tasks-to-dynamodb/fork](https://github.com/toricls/amazon-ecs-standalone-tasks-to-dynamodb/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Create a new Pull Request

## Licence

[MIT](LICENSE)

## Author

[Tori](https://github.com/toricls)
