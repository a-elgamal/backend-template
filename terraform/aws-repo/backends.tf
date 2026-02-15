terraform {
  backend "s3" {
    bucket = "myservice-tfstate"
    key    = "aws-repo/terraform.tfstate"
    region = "eu-west-1"
  }
}
