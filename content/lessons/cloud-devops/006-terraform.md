# Terraform (Infrastructure as Code)

If your Go application requires an AWS S3 Bucket, a PostgreSQL Database, and a Kubernetes cluster, you *could* log into the AWS Console and click around to create them.

However, manual clicks cannot be version-controlled, cannot be code-reviewed, and cannot be replicated easily across Dev/Staging/Prod environments.

**Terraform** solves this. It allows you to write your infrastructure as code (IaC) using HCL (HashiCorp Configuration Language). Coincidentally, Terraform itself is written entirely in Go!

## 1. The HCL Syntax

Here is an example of provisioning an AWS S3 bucket and an RDS Postgres Database for our Go application.

```hcl
# main.tf

# 1. Define the Provider (AWS)
provider "aws" {
  region = "us-east-1"
}

# 2. Provision an S3 Bucket for file uploads
resource "aws_s3_bucket" "app_uploads" {
  bucket = "my-go-app-uploads-bucket-prod"
}

# 3. Provision a Postgres Database
resource "aws_db_instance" "postgres" {
  identifier           = "go-app-db"
  allocated_storage    = 20
  engine               = "postgres"
  engine_version       = "15.3"
  instance_class       = "db.t3.micro"
  username             = "admin"
  password             = var.db_password # Injected securely via variables
  skip_final_snapshot  = true
}
```

## 2. The Execution Workflow

Terraform relies on three primary commands:

1. **`terraform init`**: Downloads the AWS provider plugins (written in Go) required to talk to the cloud API.
2. **`terraform plan`**: Connects to AWS, checks what currently exists, and prints out a "diff" (e.g., "+ aws_s3_bucket.app_uploads will be created"). **It does not make changes yet.** This allows you to review the impact before executing.
3. **`terraform apply`**: Executes the API calls to AWS to physically create the infrastructure.

## 3. The State File (`terraform.tfstate`)

How does Terraform know what already exists in AWS? 
When you run `terraform apply`, it saves a JSON file called `terraform.tfstate` locally. This file maps your HCL code to the physical AWS IDs (e.g., "my bucket maps to AWS ARN 12345").

**⚠️ THE STATE TRAP:**
If you delete your `terraform.tfstate` file, Terraform suffers amnesia. The next time you run `terraform apply`, it will think the database doesn't exist, and will attempt to create a brand new database, failing due to name collisions!

In enterprise environments, you **never** store the state file locally. You configure a "Remote Backend" to save the state file in an S3 Bucket, secured by a DynamoDB lock, so your entire team shares the same brain.

```hcl
terraform {
  backend "s3" {
    bucket         = "my-company-terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks" # Prevents two devs from applying concurrently!
  }
}
```
