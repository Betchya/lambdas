# Serverless land!

## When adding a new Lambda function

Before pushing your code, make sure you make an empty Lambda project through the AWS Console.

Create a new directory in the /lambda directory and give it the same name as the Lambda function you made through the console. You must name your new folder the same as the empty Lambda function you just created. The workflow uses the directory name to determine where to deploy the Lambda.

## If this is a new Go Lambda function

Make sure you use the provided.al2023 runtime, provided.al2 will cause runtime errors! The workflow compiles the code to be compatiable with x86-64 architecture; make sure select that when creating a new Lambda in the AWS Console.

## If this is a new NodeJS Lambda function

Select the Node.js 20.x runtime when creating a new lambda in the console.

## What does the workflow do?

When you submit a PR, a workflow will run that will detect new or updated Lambdas. On merge with main, another workflow will deploy your code to AWS.

## If credentials are failing

Go to settings -> Secrets and Variables -> Actions
Enter AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN from your AWS account
