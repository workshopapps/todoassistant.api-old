



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
                sh "sudo rm -rf /home/samuraiaj/ticked-back/*"
                sh "sudo cp -rf * /home/samuraiaj/ticked-back/"
                sh "sudo rm -rf /home/samuraiaj/ticked-back/app.env"
                sh "sudo rm -rf /home/samuraiaj/ticked-back/taskman-firebase-adminsdk.json"
                sh "sudo cp -rf /home/samuraiaj/environment/* /home/samuraiaj/ticked-back/ "
                sh "sudo systemctl restart ticked.service"
                
              
  
                
            }
        }
    }

}
