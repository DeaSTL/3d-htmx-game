pipeline {
    agent any

    parameters {
        choice(name: 'ACTION', choices: ['deploy', 'restart', 'stop', 'logs'], description: 'Action to perform')
    }

    environment {
        COMPOSE_FILE = 'docker-compose.prod.yml'
    }

    stages {
        stage('Build & Push') {
            when { expression { params.ACTION == 'deploy' } }
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-creds',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh '''
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                        docker buildx build \
                            --platform linux/amd64 \
                            -t yeahscience/htmx-3d:latest \
                            --push \
                            .
                    '''
                }
            }
        }

        stage('Deploy') {
            when { expression { params.ACTION == 'deploy' } }
            steps {
                sh '''
                    docker compose pull
                    docker compose down || true
                    docker compose up -d
                '''
            }
        }

        stage('Restart') {
            when { expression { params.ACTION == 'restart' } }
            steps {
                sh 'docker compose restart'
            }
        }

        stage('Stop') {
            when { expression { params.ACTION == 'stop' } }
            steps {
                sh 'docker compose down'
            }
        }

        stage('Logs') {
            when { expression { params.ACTION == 'logs' } }
            steps {
                sh 'docker compose logs --tail=100'
            }
        }
    }

    post {
        always {
            sh 'docker logout || true'
        }
    }
}
