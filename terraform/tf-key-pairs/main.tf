terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-1"
}

resource "aws_key_pair" "cred_laptop" {
  key_name   = "sharranm@CREDBLRMAC1270.local"
  public_key = file("/Users/sharranm/.ssh/id_rsa.pub")
}
