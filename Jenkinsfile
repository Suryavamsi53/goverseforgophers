pipeline {
    agent any

    environment {
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
                sh 'go test -v -coverprofile=coverage.out ./...'
            }
        }

        stage('Build') {
            environment {
                CGO_ENABLED = '0'
            }
            steps {
                sh 'go build -v -o bin/goverseforgophers ./cmd/server'
            }
        }

        stage('Archive') {
            steps {
                archiveArtifacts artifacts: 'bin/*', fingerprint: true
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
