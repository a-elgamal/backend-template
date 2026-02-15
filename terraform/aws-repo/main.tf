resource "aws_ecr_repository" "default" {
  name                 = "myservice"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "default" {
  repository = aws_ecr_repository.default.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Delete untagged images older than 1 day"
      selection = {
        tagStatus   = "untagged"
        countType   = "sinceImagePushed"
        countUnit   = "days"
        countNumber = 1
      }
      action = {
        type = "expire"
      }
    }]
  })
}
