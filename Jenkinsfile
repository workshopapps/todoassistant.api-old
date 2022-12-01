pipeline {

    agent any 
    tools {
      go '1.19.3'
    }

    environment {
        GO111MODULE = "on"
        CGO_ENABLED = 0
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }

    stages {
        stage("Build") {

            steps {
                
                echo "BUILD EXECUTION STARTED"
                sh "go build -o main"
            }
            
        }

        stage("Deploy") {

            steps {
                
                sh "sudo cp -rf main /home/samuraiaj/ticked-final/"
                
            }
        }
    }

}
