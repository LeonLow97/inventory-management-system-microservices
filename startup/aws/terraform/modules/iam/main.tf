data "aws_iam_policy_document" "ec2_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "api_gateway" {
  statement {
    actions = ["ssm:GetParameter"]
    resources = [
      "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/elasticache/endpoint",
      "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/auth/jwt-secret",
      "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/auth/jwt-domain",
    ]
  }

  statement {
    actions   = ["kms:Decrypt"]
    resources = ["*"]

    condition {
      test     = "StringEquals"
      variable = "kms:ViaService"
      values   = ["ssm.${var.region}.amazonaws.com"]
    }

    condition {
      test     = "StringLike"
      variable = "kms:EncryptionContext:PARAMETER_ARN"
      values   = ["arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/auth/*"]
    }
  }
}

data "aws_iam_policy_document" "authify_service" {
  statement {
    actions = ["ssm:GetParameter"]
    resources = [
      "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/rds/*",
      "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/elasticache/endpoint",
      "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/auth/jwt-secret",
      "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/s3/bucket/schema",
    ]
  }

  statement {
    actions   = ["kms:Decrypt"]
    resources = ["*"]

    condition {
      test     = "StringEquals"
      variable = "kms:ViaService"
      values   = ["ssm.${var.region}.amazonaws.com"]
    }

    condition {
      test     = "StringLike"
      variable = "kms:EncryptionContext:PARAMETER_ARN"
      values = [
        "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/rds/*",
        "arn:aws:ssm:${var.region}:${var.account_id}:parameter/authify/auth/*",
      ]
    }
  }

  statement {
    actions   = ["s3:GetObject"]
    resources = ["${var.schema_bucket_arn}/*"]
  }
}

resource "aws_iam_role" "api_gateway" {
  name               = "${var.name_prefix}-api-gateway"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume_role.json

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-api-gateway"
  })
}

resource "aws_iam_role" "authify_service" {
  name               = "${var.name_prefix}-authify-service"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume_role.json

  tags = merge(var.common_tags, {
    Name = "${var.name_prefix}-authify-service"
  })
}

resource "aws_iam_role_policy_attachment" "api_gateway_ssm" {
  role       = aws_iam_role.api_gateway.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "authify_service_ssm" {
  role       = aws_iam_role.authify_service.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "api_gateway_ecr" {
  role       = aws_iam_role.api_gateway.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_role_policy_attachment" "authify_service_ecr" {
  role       = aws_iam_role.authify_service.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_role_policy" "api_gateway" {
  name   = "${var.name_prefix}-api-gateway-runtime"
  role   = aws_iam_role.api_gateway.id
  policy = data.aws_iam_policy_document.api_gateway.json
}

resource "aws_iam_role_policy" "authify_service" {
  name   = "${var.name_prefix}-authify-service-runtime"
  role   = aws_iam_role.authify_service.id
  policy = data.aws_iam_policy_document.authify_service.json
}

resource "aws_iam_instance_profile" "api_gateway" {
  name = "${var.name_prefix}-api-gateway"
  role = aws_iam_role.api_gateway.name
}

resource "aws_iam_instance_profile" "authify_service" {
  name = "${var.name_prefix}-authify-service"
  role = aws_iam_role.authify_service.name
}
