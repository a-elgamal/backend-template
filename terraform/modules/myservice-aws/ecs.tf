resource "aws_ecs_cluster" "default" {
  name = "myservice-${var.env}"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecs_task_definition" "default" {
  family                   = "myservice-${var.env}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name      = "app-server"
      image     = var.docker_image
      essential = true

      portMappings = [{
        containerPort = var.port
        protocol      = "tcp"
      }]

      environment = [
        { name = "SERVER_HTTP_ADDRESS", value = ":${var.port}" },
        { name = "SERVER_SHUTDOWN_TIMEOUT_SECONDS", value = "30" },
        { name = "CORS_ALLOWED_ORIGINS", value = "[]" },
        { name = "DB_URL", value = "postgres://${aws_db_instance.default.username}:${random_password.db_password.result}@${aws_db_instance.default.endpoint}/${aws_db_instance.default.db_name}?sslmode=require&connect_timeout=3" },
        { name = "TELEMETRY_LOGGING_CONSOLE_LOGGING_ENABLED", value = "FALSE" },
        { name = "AWS_ALB_REGION", value = var.region },
      ]

      dependsOn = [{
        containerName = "collector"
        condition     = "HEALTHY"
      }]

      healthCheck = {
        command     = ["CMD-SHELL", "wget -qO- http://localhost:${var.port}/health || exit 1"]
        interval    = 10
        timeout     = 5
        retries     = 3
        startPeriod = 30
      }

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.default.name
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "app"
        }
      }
    },
    {
      name      = "collector"
      image     = "otel/opentelemetry-collector-contrib:0.99.0"
      essential = false

      healthCheck = {
        command     = ["CMD-SHELL", "wget -qO- http://localhost:13133/ || exit 1"]
        interval    = 10
        timeout     = 5
        retries     = 3
        startPeriod = 10
      }

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.default.name
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "collector"
        }
      }
    },
  ])
}

resource "aws_ecs_service" "default" {
  name            = "myservice"
  cluster         = aws_ecs_cluster.default.id
  task_definition = aws_ecs_task_definition.default.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [aws_security_group.ecs.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.default.arn
    container_name   = "app-server"
    container_port   = var.port
  }

  depends_on = [aws_lb_listener.https]
}
