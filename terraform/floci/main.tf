module "network" {
  source = "./modules/aws-network"

  aws_region = var.aws_region
  vpc_cidr   = var.network_cidr
}
