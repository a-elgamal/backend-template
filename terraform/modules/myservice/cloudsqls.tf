resource "google_sql_database_instance" "default" {
  name             = "myservice"
  project          = var.project
  region           = var.region
  database_version = "POSTGRES_15"
  settings {
    tier              = var.cloudsql_instance_type
    availability_type = var.cloudsql_ha_enabled ? "REGIONAL" : "ZONAL"

    ip_configuration {
      ipv4_enabled                                  = false
      private_network                               = data.google_compute_network.default.id
      enable_private_path_for_google_cloud_services = true
    }

    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }

    insights_config {
      query_insights_enabled = true
    }
  }

  deletion_protection = "true"
}

resource "google_sql_database" "default" {
  name     = "myservice"
  project  = var.project
  instance = google_sql_database_instance.default.name
}

resource "random_password" "sql_password" {
  length  = 16
  special = false
}

resource "google_sql_user" "default" {
  name     = "myservice"
  instance = google_sql_database_instance.default.name
  type     = "BUILT_IN"
  project  = var.project
  password = random_password.sql_password.result
}
