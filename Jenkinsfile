pipeline {
    agent any

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

        stage('Verify Environment') {
            steps {
                bat 'go version'
            }
        }

        stage('Dependencies') {
            steps {
                bat 'go mod download'
                bat 'go mod verify'
            }
        }

        stage('Test & Vet') {
            steps {
                bat 'go vet ./...'
                bat 'go test -v ./...'
            }
        }

        stage('Build') {
            steps {
                bat 'go build -v -o bin\\goverseforgophers.exe .'
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo 'Build succeeded!'
        }
        failure {
            echo 'Build failed. Check console output for details.'
        }
    }
}
