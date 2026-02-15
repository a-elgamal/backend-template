resource "google_cloud_run_v2_service" "default" {
  depends_on = [
    google_secret_manager_secret_iam_member.sa_secret_access,
    google_secret_manager_secret_version.otlp_collector_config,
    google_service_account_iam_member.repo_sa,
    google_project_iam_member.sa_cloudsql_user,
    google_project_iam_member.sa_cloudsql_client,
    google_project_service.sql_admin,
  ]

  name     = "myservice"
  project  = var.project
  location = var.region
  ingress  = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"


  template {
    service_account = google_service_account.default.email
    volumes {
      name = "collector-config"
      secret {
        secret = google_secret_manager_secret.otlp_collector_config.secret_id
        items {
          version = "latest"
          path    = "config.yaml"
        }
      }
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.default.connection_name]
      }
    }

    vpc_access {
      network_interfaces {
        network    = data.google_compute_network.default.id
        subnetwork = data.google_compute_subnetwork.default.id
      }
      egress = "ALL_TRAFFIC"
    }

    containers {
      name       = "app-server"
      image      = var.docker_image_digest != null ? "${var.repo_region}-docker.pkg.dev/knz-myservice-repo/knz-myservice/myservice@${var.docker_image_digest}" : "${var.repo_region}-docker.pkg.dev/knz-myservice-repo/knz-myservice/myservice:${var.docker_image_tag}"
      depends_on = ["collector"]

      resources {
        cpu_idle = false
      }

      ports {
        container_port = var.port
      }

      env {
        name  = "SERVER_HTTP_ADDRESS"
        value = ":${var.port}"
      }

      env {
        name  = "CORS_ALLOWED_ORIGINS"
        value = "[]"
      }

      env {
        name  = "DB_URL"
        value = "postgres:///${google_sql_database.default.name}?host=/cloudsql/${google_sql_database_instance.default.connection_name}&user=${google_sql_user.default.name}&password=${random_password.sql_password.result}&connect_timeout=3"
      }

      env {
        name  = "TELEMETRY_LOGGING_CONSOLE_LOGGING_ENABLED"
        value = "FALSE"
      }

      env {
        name  = "GCP_PROJECT_NUMBER"
        value = data.google_project.default.number
      }

      env {
        name  = "GCP_REGION"
        value = var.region
      }

      env {
        name  = "GCP_INTERNAL_BACKEND_SERVICE_ID"
        value = google_compute_region_backend_service.internal.generated_id
      }

      startup_probe {
        period_seconds    = 1
        failure_threshold = 60
        tcp_socket {
        }
      }
      liveness_probe {
        http_get {
          path = "/health"
        }
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }
    }

    containers {
      name  = "collector"
      image = "otel/opentelemetry-collector-contrib:0.99.0"

      startup_probe {
        http_get {
          port = 13133
          path = "/"
        }
      }

      liveness_probe {
        http_get {
          port = 13133
          path = "/"
        }
      }

      volume_mounts {
        name       = "collector-config"
        mount_path = "/etc/otelcol-contrib/"
      }
    }
  }
}

resource "google_cloud_run_service_iam_binding" "myservice_public" {
  project  = var.project
  location = google_cloud_run_v2_service.default.location
  service  = google_cloud_run_v2_service.default.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}
