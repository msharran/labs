module "network" {
  source = "./modules/aws-network"

  vpc_cidr = var.network_cidr
}
