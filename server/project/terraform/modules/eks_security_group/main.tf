resource "aws_security_group" "eks_worker_sg" {
    name = "EKS-Worker-SG"
    description = "Security group for EKS worker nodes"
    vpc_id = var.vpc_id

    # Allow inbound communication from the EKS control plane (from the Kubernetes API)
    ingress {
        from_port = 443
        to_port = 443
        protocol = "tcp"
        cidr_blocks = ["10.0.0.0/16"]
    }
    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = [ "10.0.0.0/16" ]
    }

    # Allow worker nodes to communicate with other worker nodes in the VPC (for pod communication)
    ingress {
        from_port = 10250
        to_port = 10250
        protocol = "tcp"
        cidr_blocks = [ "10.0.0.0/16" ] # allow communication within VPC
    }

    # Allow worker nodes to access microservices
    ingress {
        from_port = 443
        to_port = 443
        protocol = "tcp"
        security_groups = [ 
            var.api_gateway_sg.id,
            var.auth_sg.id
        ]
    }

    # Allow outbound traffic to other services or internet via NAT Gateway
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = [ "0.0.0.0/0" ]
    }
}
