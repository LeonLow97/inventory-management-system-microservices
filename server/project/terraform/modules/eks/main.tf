resource "aws_eks_cluster" "main" {
    name = "ims-eks-cluster"
    role_arn = ""

    vpc_config {
      subnet_ids = var.private_subnets
    }
}
