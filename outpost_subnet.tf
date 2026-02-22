resource "aws_subnet" "outpost_subnet" {
  vpc_id            = aws_vpc.outpost_vpc.id
  cidr_block        = "10.0.1.0/24"
  # Replace with your actual Outpost ARN
  outpost_arn       = "arn:aws:outposts:us-west-2:123456789012:outpost/op-0123456789abcdefg"
  availability_zone = "us-west-2a" # Must be the parent AZ of the Outpost

  tags = {
    Name = "eks-outpost-private-subnet"
  }
}