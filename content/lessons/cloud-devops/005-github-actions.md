# GitHub Actions (CI/CD)

Continuous Integration and Continuous Deployment (CI/CD) is the heartbeat of a DevOps pipeline. Every time a developer pushes Go code, an automated server should build it, test it, and deploy it.

GitHub Actions has become the industry standard for CI/CD due to its deep integration with the codebase.

## 1. The CI Pipeline (Pull Requests)

When a developer opens a Pull Request, you want to guarantee the code compiles, the tests pass, and the code is formatted correctly.

Create a file at `.github/workflows/ci.yml`:

```yaml
name: Go CI

# Trigger this workflow on all Pull Requests
on:
  pull_request:
    branches: [ "main" ]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
    # 1. Checkout the source code
    - name: Checkout code
      uses: actions/checkout@v4

    # 2. Setup the Go environment
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'
        cache: true # HUGE SPEEDUP! Automatically caches downloaded Go Modules

    # 3. Check code formatting
    - name: Verify go fmt
      run: if [ "$(go fmt ./...)" != "" ]; then echo "Code is not formatted!"; exit 1; fi

    # 4. Run the Go Linter
    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v4
      with:
        version: latest

    # 5. Run the Test Suite with the Race Detector
    - name: Run tests
      run: go test -v -race -cover ./...
```

## 2. The CD Pipeline (Deployment)

When the Pull Request is merged into the `main` branch, a CD pipeline should trigger to build the Docker image and deploy it.

```yaml
name: Go CD

# Trigger ONLY when code is merged into main
on:
  push:
    branches: [ "main" ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    # 1. Login to a Container Registry (e.g., Docker Hub or AWS ECR)
    - name: Login to Docker Hub
      uses: docker/login-action@v3
      with:
        username: ${{ secrets.DOCKERHUB_USERNAME }}
        password: ${{ secrets.DOCKERHUB_TOKEN }}

    # 2. Build and push the 10MB Scratch image
    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: .
        push: true
        # Use the git commit SHA as the image tag!
        tags: myusername/go-api:${{ github.sha }}

    # 3. Trigger Kubernetes Deployment (Example using Helm)
    # - name: Deploy to K8s
    #   run: |
    #     helm upgrade --install api ./chart \
    #       --set image.tag=${{ github.sha }}
```

## 3. The Matrix Strategy

If you are building an open-source Go CLI tool (like Terraform) that users will run on their laptops, you need to compile binaries for Linux, Mac, and Windows.

GitHub Actions allows you to use a **Matrix Strategy** to run the same job across multiple operating systems in parallel!

```yaml
jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    - name: Build binary
      run: go build -o mycli main.go
```
