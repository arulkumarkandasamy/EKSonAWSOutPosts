module "eks_outposts" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name    = var.cluster_name
  cluster_version = "1.31"

  # Outpost specific: Use private subnets on the Outpost hardware
  vpc_id     = aws_vpc.outpost_vpc.id
  subnet_ids = [aws_subnet.outpost_subnet.id]

  # Outpost nodes must be self-managed or specific EC2 instances
  self_managed_node_groups = {
    outpost_nodes = {
      instance_type = "m5.large" # Must match Outpost capacity
      min_size     = 3
      max_size     = 10
      
      # Crucial: Ensure nodes join the Outpost-based control plane
      target_group_arns = [aws_lb_target_group.outpost_tg.arn]
    }
  }

  # Enable IRSA for the Go Operator to talk to AWS APIs
  enable_irsa = true 
}