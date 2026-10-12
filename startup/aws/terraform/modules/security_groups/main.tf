resource "aws_security_group" "alb" {
  name        = "${var.name_prefix}-alb"
  description = "Allow public HTTP traffic to the Authify API."
  vpc_id      = var.vpc_id

  ingress {
    description = "Public HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-alb"
  })
}

resource "aws_security_group" "api_gateway" {
  name        = "${var.name_prefix}-api-gateway"
  description = "Allow HTTP from the load balancer."
  vpc_id      = var.vpc_id

  ingress {
    description     = "HTTP from the ALB"
    from_port       = var.api_gateway_port
    to_port         = var.api_gateway_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-api-gateway"
  })
}

resource "aws_security_group" "authify_service" {
  name        = "${var.name_prefix}-authify-service"
  description = "Allow gRPC only from API Gateway."
  vpc_id      = var.vpc_id

  ingress {
    description     = "gRPC from API Gateway"
    from_port       = var.authify_service_port
    to_port         = var.authify_service_port
    protocol        = "tcp"
    security_groups = [aws_security_group.api_gateway.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-authify-service"
  })
}

resource "aws_security_group" "postgres" {
  name        = "${var.name_prefix}-postgres"
  description = "Allow PostgreSQL only from Authify Service."
  vpc_id      = var.vpc_id

  ingress {
    description     = "PostgreSQL from Authify Service"
    from_port       = var.postgres_port
    to_port         = var.postgres_port
    protocol        = "tcp"
    security_groups = [aws_security_group.authify_service.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-postgres"
  })
}

resource "aws_security_group" "redis" {
  name        = "${var.name_prefix}-redis"
  description = "Allow Redis only from Authify application instances."
  vpc_id      = var.vpc_id

  ingress {
    description     = "Redis from API Gateway"
    from_port       = var.redis_port
    to_port         = var.redis_port
    protocol        = "tcp"
    security_groups = [aws_security_group.api_gateway.id]
  }

  ingress {
    description     = "Redis from Authify Service"
    from_port       = var.redis_port
    to_port         = var.redis_port
    protocol        = "tcp"
    security_groups = [aws_security_group.authify_service.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-redis"
  })
}
