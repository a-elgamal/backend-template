locals {
  region             = "eu-west-1"
  env                = "dev"
  vpc_id             = "TODO_VPC_ID"
  public_subnet_ids  = ["TODO_PUBLIC_SUBNET_1", "TODO_PUBLIC_SUBNET_2"]
  private_subnet_ids = ["TODO_PRIVATE_SUBNET_1", "TODO_PRIVATE_SUBNET_2"]
  certificate_arn    = "TODO_ACM_CERTIFICATE_ARN"
  ecr_account_id     = "TODO_AWS_ACCOUNT_ID"
}
