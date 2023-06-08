#!groovy

pipeline {
    agent any

    options {
        disableConcurrentBuilds()
        preserveStashes(buildCount: 7)
        buildDiscarder(logRotator(numToKeepStr: '150', artifactNumToKeepStr: '150'))
    }
    triggers {
        pollSCM('H/3 * * * *')
    }
    stages {
        stage('prepare') {
            agent {
                docker {
                    image 'ubuntu'
                    alwaysPull true
                }
            }
            steps {
                sh 'env'
            }
        }
}
