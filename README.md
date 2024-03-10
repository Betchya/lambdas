# Serverless Land: Deployment Guide

Welcome to the guide for deploying Lambda functions. Here's what you need to know:

## Workflow Overview

### 1. Generate Lambda List File

Upon creating a pull request targeting the main branch, this workflow:

- Identifies modified Lambda functions.
- Generates a list of these functions and uploads it as an artifact.
- Copies the list to an S3 bucket for tracking changes.

### 2. Deploy Lambda Functions

When changes are merged into the main branch, this workflow:

- Sets up the required environments for Go and Node.js.
- Downloads the list of modified Lambda functions from S3.
- Builds, zips, and deploys each function to AWS Lambda.

## Adding a New Lambda Function

### AWS Console Setup

1. **Create an Empty Lambda Function**: Start by creating a new Lambda function through the AWS Console. This initial step is crucial for setting up the infrastructure on AWS.

2. **Directory Naming**: In your project, create a new directory under `/lambdas` with the **exact name** of your newly created Lambda function. This naming consistency is vital for our workflows to recognize and deploy your function correctly.

### Go Lambda Functions

- **Runtime Requirement**: Utilize the `provided.al2023` runtime. Avoid `provided.al2` to prevent runtime errors.
- **Architecture**: Ensure your function is compiled for `x86-64` architecture, matching your AWS Lambda creation settings.

### Node.js Lambda Functions

- **Runtime Requirement**: Opt for the `Node.js 20.x` runtime when creating your Lambda function in the AWS Console.
- **Include a `package.json`**: Your function's directory must contain a `package.json` file to manage dependencies and scripts.

## Deployment Process

- **Pull Requests**: Submitting a PR triggers the "Generate Lambda List File" workflow, identifying any Lambda function changes.
- **Merging**: On merging a PR to the main branch, the "Deploy Lambda Functions" workflow activates, deploying your updates to AWS Lambda.

### Special Notes

- **Build and Test Scripts**: For Node.js Lambda functions, ensure your `package.json` defines `build` and `test` scripts. Our workflow will execute these if present, preparing your function for deployment.

## Troubleshooting AWS Credentials

If deployment fails due to AWS credential issues:

1. Go to **Settings** > **Secrets and Variables** > **Actions** in your GitHub repository.
2. Ensure `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and `AWS_SESSION_TOKEN` are correctly set with your AWS account details.
