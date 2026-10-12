# AWS EC2 with Terraform

This guide creates the Authify stack in AWS with Terraform. It uses two Availability Zones, Spot application instances, one NAT Gateway, RDS PostgreSQL, Redis, ECR, and an Application Load Balancer.

## Prerequisites

- Terraform 1.5 or newer
- AWS CLI configured for the account and region you want to use
- Docker
- `curl`
- AWS permissions to manage VPC, EC2, RDS, ElastiCache, ECR, S3, IAM, SSM, and load balancing resources, plus read access to the selected Availability Zones and public Amazon Linux AMI parameter

Terraform creates and deletes the AWS infrastructure. The AWS CLI provides credentials and handles ECR login, image publishing, and runtime access to SSM parameters and ECR from the EC2 instances.

If an older AWS CLI stack is still running in this account, remove it with `make -C startup/aws/cli clean` before creating the Terraform stack. Terraform starts fresh resources and does not adopt the CLI-created resources.

## Set up secrets

From the repository root, copy the example variables file:

```sh
cp startup/aws/terraform/environments/ec2/terraform.tfvars.example \
  startup/aws/terraform/environments/ec2/terraform.tfvars
```

Set strong values for `master_password`, `authify_password`, and `jwt_secret` in `terraform.tfvars`. Keep this file and Terraform state private. The repository ignores both.

Terraform uses your standard AWS CLI profile and defaults to `ap-southeast-1`. To use another profile, set `AWS_PROFILE` before running the Make commands.

The Make targets supply `schema_file` automatically from the repository root. To use another schema, pass `SCHEMA_FILE=/absolute/path/users.sql` to Make.

For direct Terraform `plan`, `apply`, or `destroy` commands, set the schema input first from the repository root:

```sh
export TF_VAR_schema_file="$(pwd)/authify-service/schema/users.sql"
```

## Create the stack

Preview the Terraform changes:

```sh
make -C startup aws-ec2-plan
```

Create the infrastructure, build and push versioned container images, then wait for both API Gateway instances to pass the readiness check:

```sh
make -C startup aws-ec2-up
```

The command prints the API URL after the load balancer and both application instances are ready. You can retrieve it later with:

```sh
terraform -chdir=startup/aws/terraform/environments/ec2 output -raw api_base_url
```

Check readiness with:

```sh
curl "$(terraform -chdir=startup/aws/terraform/environments/ec2 output -raw api_base_url)/health/readiness"
```

To use the same image tag for a plan and apply, pass `IMAGE_TAG` to both commands, for example `IMAGE_TAG=portfolio-01`.
Use a new tag for each deployment because the ECR repository treats image tags as immutable.

## Delete the stack

This portfolio setup is intended to be removed after use. Destroy the Terraform-managed resources with:

```sh
make -C startup aws-ec2-down
```

Keep `terraform.tfvars` until destroy completes. The database skips its final snapshot, and Terraform deletes the ECR images and schema bucket contents.

## Design notes

- Spot instances reduce EC2 cost and can be interrupted.
- A single NAT Gateway reduces cost. Both private application subnets use it, so outbound traffic from the second Availability Zone crosses zones.
- Database and Redis are private and accept traffic only from the application security groups.
- EC2 access uses Systems Manager Session Manager. The stack does not create public SSH hosts or key pairs.
- PostgreSQL and Redis use single-node, non-Multi-AZ configurations to keep this portfolio stack small.
- The public load balancer serves HTTP only. Use throwaway credentials for portfolio testing. HTTPS would require a domain and an ACM certificate.
- Terraform state is local and contains sensitive values from the variables file. Do not commit or share it. For a team or long-lived environment, move state to a private encrypted backend with locking.

The previous AWS CLI provisioning scripts remain under `startup/aws/cli` for reference and cleanup of older CLI-created stacks. Use this Terraform guide for new stacks.
