resource "aws_ecr_repository" "this" {
  name                 = var.ecr_repository_name
  image_tag_mutability = "IMMUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-images"
  })
}

resource "aws_instance" "authify_service" {
  for_each = var.app_subnets

  ami                         = var.ami_id
  instance_type               = var.authify_service_instance_type
  subnet_id                   = each.value
  vpc_security_group_ids      = [var.authify_service_security_group_id]
  iam_instance_profile        = var.authify_service_instance_profile
  associate_public_ip_address = false
  user_data_replace_on_change = true
  user_data = templatefile("${path.module}/user_data/authify-service.sh.tftpl", {
    region              = var.region
    image_uri           = "${aws_ecr_repository.this.repository_url}:authify-service-${var.image_tag}"
    ecr_registry        = split("/", aws_ecr_repository.this.repository_url)[0]
    service_port        = var.authify_service_port
    postgres_port       = var.postgres_port
    redis_port          = var.redis_port
    parameter_names     = var.parameter_names
    schema_bucket_name  = var.schema_bucket_name
    schema_object_key   = var.schema_object_key
    initialize_database = each.key == "az1"
  })

  metadata_options {
    http_tokens = "required"
  }

  root_block_device {
    encrypted   = true
    volume_type = "gp3"
    volume_size = 8
  }

  dynamic "instance_market_options" {
    for_each = var.use_spot_instances ? [true] : []

    content {
      market_type = "spot"

      spot_options {
        instance_interruption_behavior = "terminate"
      }
    }
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-authify-service-${each.key}"
    Role = "authify-service"
  })
}

resource "aws_instance" "api_gateway" {
  for_each = var.app_subnets

  ami                         = var.ami_id
  instance_type               = var.api_gateway_instance_type
  subnet_id                   = each.value
  vpc_security_group_ids      = [var.api_gateway_security_group_id]
  iam_instance_profile        = var.api_gateway_instance_profile
  associate_public_ip_address = false
  user_data_replace_on_change = true
  user_data = templatefile("${path.module}/user_data/api-gateway.sh.tftpl", {
    region                  = var.region
    image_uri               = "${aws_ecr_repository.this.repository_url}:api-gateway-${var.image_tag}"
    ecr_registry            = split("/", aws_ecr_repository.this.repository_url)[0]
    authify_service_address = aws_instance.authify_service[each.key].private_ip
    authify_service_port    = var.authify_service_port
    api_gateway_port        = var.api_gateway_port
    redis_port              = var.redis_port
    parameter_names         = var.parameter_names
  })

  metadata_options {
    http_tokens = "required"
  }

  root_block_device {
    encrypted   = true
    volume_type = "gp3"
    volume_size = 8
  }

  dynamic "instance_market_options" {
    for_each = var.use_spot_instances ? [true] : []

    content {
      market_type = "spot"

      spot_options {
        instance_interruption_behavior = "terminate"
      }
    }
  }

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-api-gateway-${each.key}"
    Role = "api-gateway"
  })
}
