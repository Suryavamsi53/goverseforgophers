# Terraform (Infrastructure as Code)

## 1. Learning Objectives
* **What you'll learn**: How to provision and manage cloud infrastructure (AWS/GCP servers, databases, networks) using HashiCorp Terraform and HCL (HashiCorp Configuration Language).
* **Why it matters**: Clicking buttons in the AWS Console is dangerous. If you accidentally delete a database, you don't know how it was configured to rebuild it. Terraform treats servers exactly like Go code: version-controlled, testable, and strictly reproducible.
* **Where it's used**: The undisputed industry standard for multi-cloud infrastructure provisioning.

---

## 2. Real-world Story
Imagine building a complex LEGO castle. You spend 5 hours clicking pieces together (The AWS Console). Suddenly, your little brother drops it. It shatters. You have no idea how to rebuild it exactly as it was.
Terraform is the LEGO Instruction Manual. You write down exactly what pieces you want: "1 Red Block, 2 Blue Blocks". You hand the manual to a robot (The Terraform CLI). The robot perfectly builds the castle. If it shatters, you just hand the manual to the robot again, and it rebuilds it flawlessly in seconds.

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[main.tf (HCL Code)] -->|terraform plan| B(Terraform Core)
    B -->|Calculates Diff| C{State File: terraform.tfstate}
    
    B -->|terraform apply| D[AWS Provider]
    
    D -->|API Call| E[Create EC2 Server]
    D -->|API Call| F[(Create Postgres RDS)]
    
    style B fill:#8b5cf6,color:#fff
    style C fill:#ef4444,color:#fff
```

---

## 4. Internal Working (Under the Hood)
Terraform operates on a **Declarative State Machine**.
1. **The Code**: You declare what you want (e.g., "I want 3 Servers").
2. **The State (`tfstate`)**: A JSON file where Terraform remembers what actually exists in AWS right now (e.g., "There are currently 2 Servers").
3. **The Plan**: Terraform calculates the mathematical difference (Diff) between your Code (3) and the State (2). It determines: "I need to create 1 Server."
4. **The Apply**: It makes the API calls to AWS to match reality with your code, and updates the State file.

---

## 5. Compiler Behavior
* **HCL (HashiCorp Configuration Language)**: Terraform doesn't use JSON or YAML. It uses HCL, a custom declarative language written in Go! It strongly resembles Go structs. The Terraform CLI compiles this HCL into an internal Graph data structure to execute dependencies in parallel.

---

## 6. Memory Management
* **The Provider Architecture**: The Terraform Core engine doesn't know how to talk to AWS. It dynamically downloads a Go binary called a "Provider" (e.g., `terraform-provider-aws`). The Core engine communicates with the Provider over local gRPC, ensuring massive cloud configurations don't bloat the core memory footprint.

---

## 7. Code Examples

### 🔹 Example 1: Simple (Provisioning a Server)
```hcl
# main.tf
provider "aws" {
  region = "us-east-1"
}

# Resource Type, Resource Name
resource "aws_instance" "go_server" {
  ami           = "ami-0c55b159cbfafe1f0" # Ubuntu Linux
  instance_type = "t3.micro"
  
  tags = {
    Name = "GoVerse-API"
  }
}
```

### 🔹 Example 2: Intermediate (Variables and Outputs)
```hcl
variable "server_count" {
  type    = number
  default = 3
}

resource "aws_instance" "cluster" {
  count         = var.server_count # Provisions 3 identical servers!
  ami           = "ami-xyz"
  instance_type = "t3.micro"
}

output "server_ips" {
  # Prints the public IPs of all 3 servers to the terminal!
  value = aws_instance.cluster[*].public_ip 
}
```

### 🔹 Example 3: Advanced (Dependency Graph)
```hcl
# Creating a Database and a Server, passing the DB URL into the Server!
resource "aws_db_instance" "postgres" {
  engine = "postgres"
  # ... other config
}

resource "aws_instance" "go_app" {
  ami = "ami-xyz"
  
  # Implicit Dependency! Terraform mathematically realizes it MUST wait 
  # for the DB to be fully created before booting the server!
  user_data = "export DB_URL=${aws_db_instance.postgres.endpoint}"
}
```

### 🔹 Example 4: Production (Remote State Storage)
```hcl
# NEVER store terraform.tfstate on your local laptop in Git!
# Store it in a centralized AWS S3 bucket with DynamoDB locking.
terraform {
  backend "s3" {
    bucket         = "goverse-terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "terraform-locks" # Prevents two devs from applying at once!
  }
}
```

### 🔹 Example 5: Interview
```hcl
# Q: If I manually delete the AWS server in the AWS Web Console, what happens when I run `terraform plan`?
# A: Terraform compares the State with reality, realizes the server is missing, 
# and the Plan will output: "1 to create". It will instantly rebuild the missing server!
```

---

## 8. Production Examples
1. **Multi-Cloud**: Using Terraform to provision an AWS EC2 instance, a Cloudflare DNS record, and a DataDog monitoring dashboard in the exact same `main.tf` file.
2. **Modules**: Creating a custom Terraform Module (a reusable block of HCL). Instead of defining a complex VPC (Network) in 500 lines of code, you just call `module "network" { source = "./vpc" }`.

---

## 9. Performance & Benchmarking
* **Parallel API Execution**: Because Terraform builds a Directed Acyclic Graph (DAG) of dependencies, if you declare 10 AWS S3 buckets and 10 Google Cloud SQL databases, it will fire 20 asynchronous Goroutines and provision them all simultaneously in 3 seconds!

---

## 10. Best Practices
* ✅ **Do**: ALWAYS run `terraform plan` and read the output carefully before running `terraform apply`.
* ❌ **Don't**: Modify resources manually in the AWS Web Console. This causes "State Drift". Terraform will aggressively overwrite your manual changes the next time it runs to force reality back to what the HCL code dictates!
* 🏢 **Google / Uber / Netflix Style**: Use Atlantis. It is a Go application that runs inside your GitHub repo. When you make a Pull Request altering `main.tf`, Atlantis automatically runs `terraform plan` and pastes the output as a GitHub comment!

---

## 11. Common Mistakes
1. **Committing State to Git**: The `terraform.tfstate` file contains raw JSON of your entire infrastructure, including raw unencrypted Database Passwords! If you push this to GitHub, you will be hacked in 5 minutes. Always add `*.tfstate` to `.gitignore` and use a Remote Backend (S3).
2. **Missing `depends_on`**: If you create a Server and a Security Group, but forget to link them via variables, Terraform will try to build them in parallel. The Server will fail because the Security Group doesn't exist yet!

---

## 12. Debugging
How to troubleshoot Terraform in production:
* **Trace Logs**: Run `TF_LOG=TRACE terraform apply`. This exposes the underlying Go HTTP requests Terraform is making to the AWS API, allowing you to debug exactly why AWS rejected your configuration.

---

## 13. Exercises
1. **Easy**: Write HCL to provision an `aws_s3_bucket`.
2. **Medium**: Run `terraform init` to download the AWS provider, then `terraform plan` to see the creation logic.
3. **Hard**: Define a variable for the bucket name, and pass it in using a `terraform.tfvars` file.
4. **Expert**: Configure a Remote State Backend using an existing S3 bucket.

---

## 14. Quiz
1. **MCQ**: What is the purpose of the `terraform.tfstate` file?
   * (A) To store HCL code (B) To cache AWS credentials (C) To map real-world infrastructure IDs to your HCL code. *(Answer: C)*
2. **System Design Follow-up**: Why is Terraform called "Declarative" rather than "Imperative" (like Bash scripts)? *(In Bash, you say "Create server X. Create server X". If you run it twice, you get 2 servers. In Terraform, you say "I want exactly 1 server X". If you run it 100 times, it does nothing after the first time, because the state already matches the declaration).*

---

## 15. FAANG Interview Questions
* **Beginner**: What is Infrastructure as Code (IaC)?
* **Intermediate**: Explain the mechanism of State Locking and why it is critical for CI/CD.
* **Senior (Google/Meta)**: Explain how you would refactor a monolithic `main.tf` file managing 1,000 resources into granular Workspaces and Modules to limit the blast radius of a bad `terraform apply`.

---

## 16. Mini Project
**The Go Cloud Architecture**
* Write a Terraform script that provisions:
  1. An AWS ECR (Docker Registry)
  2. An AWS RDS (PostgreSQL Database)
  3. An AWS ECS Cluster (To run the Go Docker container)
* Ensure the ECS cluster cannot boot until the Database is fully provisioned.
* Destroy the entire infrastructure perfectly using `terraform destroy`.

---

## 17. Enterprise Features & Observability
* **Drift Detection**: Enterprise tools periodically run `terraform plan` in the background. If the plan outputs anything other than "No changes", it means a rogue engineer logged into the AWS console and manually altered production. It fires a critical PagerDuty alert!

---

## 18. Source Code Reading
Walkthrough of `hashicorp/terraform`.
* **The Graph Builder**: Look at the `terraform/graph_builder.go` file. It is a masterclass in Go computer science, traversing the Abstract Syntax Tree of the HCL to build a strict dependency graph to safely execute operations.

---

## 19. Architecture
* **Pulumi vs Terraform**: Terraform uses a custom language (HCL). Pulumi allows you to define infrastructure using actual Go code (`app := ecs.NewCluster(ctx, "app")`). Some teams prefer Pulumi because they can use standard Go `for` loops and `if` statements instead of learning HCL!

---

## 20. Summary & Cheat Sheet
* **Init**: `terraform init` (Downloads providers).
* **Plan**: `terraform plan` (Dry run).
* **Apply**: `terraform apply` (Executes).
* **State**: The mapping of Code to Reality. Never store in Git.
