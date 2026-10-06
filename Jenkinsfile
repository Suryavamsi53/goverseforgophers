pipeline {
    agent any

    tools {
        // If configured under Manage Jenkins -> Tools -> Go
        // Otherwise, ensure 'go' is installed on the Jenkins host/agent
        go 'go-latest' 
    }

    environment {
        CGO_ENABLED = '0'
        GO111MODULE = 'on'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Dependencies') {
            steps {
                sh 'go mod download'
                sh 'go mod verify'
            }
        }

        stage('Test & Vet') {
            steps {
                sh 'go vet ./...'
                sh 'go test -v -race -coverprofile=coverage.out ./...'
            }
        }

        stage('Build') {
            steps {
                sh 'go build -v -o bin/goverseforgophers .'
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo 'Build, tests, and vet checks passed successfully!'
        }
        failure {
            echo 'Build failed. Check console output for errors.'
        }
    }
}