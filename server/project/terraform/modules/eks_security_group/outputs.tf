output "eks_sg_id" {
    value = aws_security_group.eks_worker_sg.id
}