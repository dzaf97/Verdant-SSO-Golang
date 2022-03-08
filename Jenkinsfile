// node {
//     checkout scm

//     stage('Build Images') {
//             steps {
//                  echo "BUILD_IMAGES"
//                 def testimage = docker.build("test:v4")
//                 // script {
//                 //     def testimage = docker.build("test:v3")
//                 // }
//             }
//         }
//         stage('Test') {
//             steps {
//                 // sh './jenkins/scripts/test.sh'
//                 echo "TEST_PHASE"
//             }
//         }
//         stage('Deliver for development') {
//             when {
//                 branch 'development' 
//             }
//             steps {
//                 echo "DEPLOYED_STAGING"
//                 // sh './jenkins/scripts/deliver-for-development.sh'
//                 // input message: 'Finished using the web site? (Click "Proceed" to continue)'
//                 // sh './jenkins/scripts/kill.sh'
//             }
//         }
//         stage('Deploy for production') {
//             when {
//                 branch 'production'  
//             }
//             steps {
//                 echo "DEPLOYED_PROD"
//                 // sh './jenkins/scripts/deploy-for-production.sh'
//                 // input message: 'Finished using the web site? (Click "Proceed" to continue)'
//                 // sh './jenkins/scripts/kill.sh'
//             }
//         }
// }

pipeline {
    agent any
    environment {
        CI = 'true'
    }
    stages {
        stage('Build Images') {
            steps {
                echo "BUILD IMAGERS"
                script {
                    
                    testimage = docker.build("test:v3")
                }
            }
        }
        stage('Test') {
            steps {
                // sh './jenkins/scripts/test.sh'
                echo "TEST_PHASE"
            }
        }
        stage('Deliver for development') {
            when {
                branch 'development' 
            }
            steps {
                echo "DEPLOYED_STAGING"
                // sh './jenkins/scripts/deliver-for-development.sh'
                // input message: 'Finished using the web site? (Click "Proceed" to continue)'
                // sh './jenkins/scripts/kill.sh'
            }
        }
        stage('Deploy for production') {
            when {
                branch 'production'  
            }
            steps {
                echo "DEPLOYED_PROD"
                // sh './jenkins/scripts/deploy-for-production.sh'
                // input message: 'Finished using the web site? (Click "Proceed" to continue)'
                // sh './jenkins/scripts/kill.sh'
            }
        }
    }
}