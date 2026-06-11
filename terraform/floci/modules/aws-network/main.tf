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

variable "aws_region" {
  type        = string
  description = "aws region"
  default     = "us-east-1"
}

variable "public_subnets" {
  type = map(object({
    cidr = string
    az   = string
  }))
  description = "public subnets"
  default     = null
}

variable "private_subnets" {
  type = map(object({
    cidr = string
    az   = string
  }))
  description = "private subnets"
  default     = null
}

locals {
  public_subnets = coalesce(var.public_subnets, {
    public_subnet_a = {
      cidr = "10.0.0.0/22"
      az   = "${var.aws_region}a"
    },
    public_subnet_b = {
      cidr = "10.0.4.0/22"
      az   = "${var.aws_region}b"
    }
  })

  private_subnets = coalesce(var.private_subnets, {
    private_subnet_a = {
      cidr = "10.0.8.0/22"
      az   = "${var.aws_region}a"
    },
    private_subnet_b = {
      cidr = "10.0.12.0/22"
      az   = "${var.aws_region}b"
    }
  })

  public_subnet_by_az = {
    for name, subnet in local.public_subnets : subnet.az => name
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
  for_each = local.public_subnets

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value.cidr
  availability_zone = each.value.az

  tags = {
    Name = each.key
  }
}

resource "aws_subnet" "private" {
  for_each = local.private_subnets

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value.cidr
  availability_zone = each.value.az

  tags = {
    Name = each.key
  }
}

resource "aws_eip" "nat" {
  for_each = local.public_subnet_by_az

  domain = "vpc"
}

resource "aws_nat_gateway" "nat" {
  for_each = local.public_subnet_by_az

  allocation_id = aws_eip.nat[each.key].id
  subnet_id     = aws_subnet.public[each.value].id

  depends_on = [aws_internet_gateway.main]

  tags = {
    Name = "nat-${each.key}"
  }
}

resource "aws_route_table" "private" {
  for_each = local.private_subnets

  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat[each.value.az].id
  }

  tags = {
    Name = "${each.key}-private"
  }
}

resource "aws_route_table_association" "private" {
  for_each = aws_subnet.private

  subnet_id      = each.value.id
  route_table_id = aws_route_table.private[each.key].id
}
