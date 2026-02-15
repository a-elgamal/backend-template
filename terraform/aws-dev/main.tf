module "myservice" {
  source             = "../modules/myservice-aws"
  env                = local.env
  region             = local.region
  vpc_id             = local.vpc_id
  public_subnet_ids  = local.public_subnet_ids
  private_subnet_ids = local.private_subnet_ids
  certificate_arn    = local.certificate_arn
  docker_image       = var.docker_image != null ? var.docker_image : "${local.ecr_account_id}.dkr.ecr.${local.region}.amazonaws.com/myservice:latest"
  oidc_client_id     = "TODO_GOOGLE_OAUTH_CLIENT_ID"
  oidc_client_secret = var.oidc_client_secret
  rds_multi_az       = false
}
