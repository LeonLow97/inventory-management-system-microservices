resource "aws_lb" "this" {
  name               = "${var.name_prefix}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.security_group_id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = false

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-alb"
  })
}

resource "aws_lb_target_group" "api_gateway" {
  name        = "${var.name_prefix}-api-gateway"
  port        = var.api_gateway_port
  protocol    = "HTTP"
  target_type = "instance"
  vpc_id      = var.vpc_id

  health_check {
    enabled             = true
    path                = "/health/readiness"
    port                = "traffic-port"
    protocol            = "HTTP"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
    matcher             = "200"
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-api-gateway"
  })
}

resource "aws_lb_target_group_attachment" "api_gateway" {
  for_each = var.api_gateway_instance_ids

  target_group_arn = aws_lb_target_group.api_gateway.arn
  target_id        = each.value
  port             = var.api_gateway_port
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.this.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api_gateway.arn
  }
}
