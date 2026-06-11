# terraform {
#   required_providers {
#     aws = {
#       source  = "hashicorp/aws"
#       version = "~> 5.0"
#     }
#   }
# }

variable "vpc_cidr" {
  type        = string
  description = "vpc cidr block"
}

variable "public_subnets" {
  type = map(object({
    cidr              = string
    availability_zone = optional(string)
  }))
  description = "public subnets"
  default = {
    public_subnet_a = {
      cidr = "10.0.0.0/22"
    },
    public_subnet_b = {
      cidr = "10.0.4.0/22"
    }
  }
}

variable "private_subnets" {
  type = map(object({
    cidr              = string
    availability_zone = optional(string)
  }))
  description = "private subnets"
  default = {
    private_subnet_a = {
      cidr = "10.0.8.0/22"
    },
    private_subnet_b = {
      cidr = "10.0.12.0/22"
    }
  }
}

resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
}

resource "aws_route_table_association" "public" {
  for_each = aws_subnet.public

  subnet_id      = each.value.id
  route_table_id = aws_route_table.public.id
}

resource "aws_subnet" "public" {
  for_each = var.public_subnets

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value.cidr
  availability_zone = each.value.availability_zone

  tags = {
    Name = each.key
  }
}

resource "aws_subnet" "private" {
  for_each = var.private_subnets

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value.cidr
  availability_zone = each.value.availability_zone

  tags = {
    Name = each.key
  }
}
