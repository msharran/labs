# Terraform + Floci AWS practice

This directory starts a local [Floci](https://floci.io/aws/) AWS emulator with Docker Compose on [OrbStack](https://orbstack.dev/) and configures AWS/Terraform env vars via direnv.

## First-time setup

Install prerequisites if needed:

```sh
brew install --cask orbstack direnv
```

Then allow the directory env:

```sh
cd terraform/floci
direnv allow
```

After `direnv allow`, entering this directory will run `docker compose up -d` and export:

- `AWS_ENDPOINT_URL=http://localhost:4566`
- `AWS_ACCESS_KEY_ID=test`
- `AWS_SECRET_ACCESS_KEY=test`
- `AWS_DEFAULT_REGION=us-east-1`
- `TF_VAR_aws_endpoint_url=http://localhost:4566`
- `TF_VAR_aws_region=us-east-1`

## Practice workflow

```sh
terraform init
terraform plan
terraform apply

aws s3 ls
aws sqs list-queues
```

Stop Floci when done:

```sh
docker compose down
```

Delete persisted emulator state:

```sh
rm -rf data
```
