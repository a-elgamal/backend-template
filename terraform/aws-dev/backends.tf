terraform {
  backend "s3" {
    bucket = "myservice-tfstate"
    key    = "aws-dev/terraform.tfstate"
    region = "eu-west-1"
  }
}
