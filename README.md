# EKS on Outposts Namespace Operator

This project implements a Kubernetes Operator using the `controller-runtime` library. It manages the lifecycle of `Namespace` resources, specifically focusing on integrating with external AWS resources in an EKS on Outposts environment.

## Features

1.  **External Resource Cleanup**:
    - When a Namespace is deleted, the operator intercepts the deletion via a Finalizer (`aws.cleanup.finalizer`).
    - It triggers a cleanup of associated AWS RDS resources before allowing the Namespace to be removed.

2.  **Drift Detection & Provisioning**:
    - Ensures that an IAM Role for Service Accounts (IRSA) exists for the namespace.
    - Requeues if provisioning fails or drift is detected.

3.  **Security Compliance**:
    - Automatically injects mandatory security labels (`enterprise-security: compliant`) to ensure namespaces meet enterprise standards.

## Project Structure

- `main.go`: Entry point for the operator. Sets up the Manager and registers the controller.
- `controllers/operator.go`: Contains the `NamespaceReconciler` logic.
- `controllers/aws_manager.go`: A wrapper for AWS SDK interactions (currently mocked).

## Getting Started

### Prerequisites
- Go 1.19+
- Access to a Kubernetes cluster

### Running
```bash
go run main.go
```