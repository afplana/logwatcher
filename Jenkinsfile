#!/usr/bin/env groovy
pipeline {
    agent any
    tools {
        go 'go1.15.3'
    }

    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }

    stages {        
        stage('Build') {
            steps {
                echo 'Compiling and building'
                sh 'pwd'
                sh 'go version'
                echo 'Start Building process...'
                sh 'cd golang-logwatcher && go build'
            }
        }

        stage ('Archive') {
            steps{
                echo 'Archive artifacts'
                archiveArtifacts 'golang-logwatcher/golang-logwatcher'
            }
        }
    }
}