module "myservice" {
  source              = "../modules/myservice"
  env                 = local.env
  project             = local.service_project
  host_project        = local.host_project
  region              = local.env_region
  docker_image_digest = var.docker_image_digest
  cloudsql_ha_enabled = false
}
