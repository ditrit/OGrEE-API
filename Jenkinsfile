//The URL below is a good place to start
// https://raw.githubusercontent.com
// /github/platform-samples/master/hooks/jenkins/jira-workflow/Jenkinsfile
pipeline {
    //agent {dockerfile true}

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                //cd /var/lib/jenkins/workspace/Job1prototypev2
                //docker build -t testingalpine .
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                //docker rm $(docker ps -aq)
                //docker run testingalpine sh -c "cd prototypev2 && go test -v ./..."
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}