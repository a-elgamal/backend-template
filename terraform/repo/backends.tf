terraform {
  backend "gcs" {
    bucket = "knz-myservice-tfstate"
    prefix = "repo"
  }
}
