resource "aws_lb" "default" {
  name               = "myservice-${var.env}"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids
}

resource "aws_lb_target_group" "default" {
  name        = "myservice-${var.env}"
  port        = var.port
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    path                = "/health"
    protocol            = "HTTP"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    interval            = 15
    timeout             = 5
  }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.default.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = var.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.default.arn
  }
}

# OIDC authentication on /internal/* routes using Google as the identity provider
resource "aws_lb_listener_rule" "internal_oidc" {
  listener_arn = aws_lb_listener.https.arn
  priority     = 100

  action {
    type = "authenticate-oidc"

    authenticate_oidc {
      authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
      token_endpoint         = "https://oauth2.googleapis.com/token"
      user_info_endpoint     = "https://openidconnect.googleapis.com/v1/userinfo"
      issuer                 = "https://accounts.google.com"
      client_id              = var.oidc_client_id
      client_secret          = var.oidc_client_secret
      scope                  = "openid email profile"
    }
  }

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.default.arn
  }

  condition {
    path_pattern {
      values = ["/internal/*"]
    }
  }
}
