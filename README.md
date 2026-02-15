# My Service

A backend service template with a Flutter admin portal.

## Prerequisites
- Go 1.26+
- PostgreSQL 18
- Flutter 3.41+ (for portal development)
- Docker & Docker Compose (optional)

## Running all checks
```shell
docker-compose -f docker-compose.local.yml up
make
```

## Running the service & portal
```shell
docker-compose -f docker-compose.local.yml up
```
The portal will be available at http://localhost:8080/internal/portal and the service is also available on the same port.

## Running the portal for local development
1. Run the portal using your IDE or Flutter run command.
1. Update SERVER.CORS_ALLOWED_ORIGINS in application.yaml to allow the URL of the portal to be accepted by the server (Flutter typically generates a random port for Flutter run)
1. Run postgres (either locally or in a container). Make sure that the postgres credentials and database are correctly set in application.yaml
1. Run the service locally using `make run`
1. Refresh the portal in your browser.

## Customizing for your project

1. Replace `alielgamal.com/myservice` with your module path in `go.mod` and all import paths
2. Replace `myservice` with your service name in `Makefile`, `Dockerfile`, `docker-compose.local.yml`, `application.yaml`, and code
3. Update the `item` example entity in `internal/item/` with your domain entities
4. Update Terraform files — search for `TODO` comments for values that need customization
5. Update GitHub Actions workflows — search for `TODO` placeholders
