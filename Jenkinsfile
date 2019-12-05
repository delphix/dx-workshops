pipeline {
    options {
        disableConcurrentBuilds()
    }
    agent any
    environment {
        ORACLE_CHANGE = false
        PG_CHANGE = false
        ANSIBLE_CERT = "/var/lib/jenkins/.ssh/ansible"
    }
    parameters {
        booleanParam(name: 'ORACLE_TESTING', defaultValue: true, description: 'Run Testing on Change')
        booleanParam(name: 'ORACLE_FORCE_TESTING', defaultValue: false, description: 'Force Testing')
        booleanParam(name: 'PG_TESTING', defaultValue: true, description: 'Run Testing on Change')
        booleanParam(name: 'PG_FORCE_TESTING', defaultValue: false, description: 'Force Testing')
    }
    stages {
        stage('Prep Environment'){
            steps {
                library "environment-lib@${BRANCH_NAME}"
                echo "Grabbing certs"
                sh "cp ${env.ANSIBLE_CERT}* certs"
                echo "Copying env file"
                sh "cp /var/lib/jenkins/.docker.env ."
                echo "cleaning up docker"
                sh "docker rmi -f tcw-oracle tcw-pg"
                // echo "Copying s3 backend file"
                // sh "cp /var/lib/jenkins/terraform_dev_backend.tf ${env.TF_DIR}/."
            }
        }
        stage('Build Workshop Docker Containers'){
            parallel{
                stage('Build tcw-oracle Container'){
                    steps{
                        echo "Building Container"
                        sh "docker-compose build  --no-cache tcw-oracle" 
                    }
                }
                stage('Build tcw-pg Container'){
                    steps{
                        echo "Building Container"
                        sh "docker-compose build --no-cache tcw-pg" 
                    }
                }
            }
        }
        stage("Packer Builds"){
            stages{
                stage('Oracle Packer Build'){
                    steps{
                        script{
                            environment.packerBuild("tcw-oracle")
                            ORACLE_CHANGE = sh (
                                script: """#!/bin/bash
                                        { set -x; } 2>/dev/null
                                        [[ -f change.ignore ]] && echo "true" || echo "false" """,
                                returnStdout: true
                            ).trim()
                            }
                            sh "echo ${ORACLE_CHANGE}"
                    }
                }
                stage('Postgres Packer Build'){
                    steps{
                        script{
                            environment.packerBuild("tcw-pg")
                            PG_CHANGE = sh (
                                script: """#!/bin/bash
                                        { set -x; } 2>/dev/null
                                        [[ -f change.ignore ]] && echo "true" || echo "false" """,
                                returnStdout: true
                            ).trim()
                            }
                            sh "echo ${PG_CHANGE}"
                    }
                }
            }
        }
        stage('Deploy and Test'){
            parallel{
                stage('Oracle'){
                    environment {
                        WORKSHOP = "tcw-oracle"
                        TF_DIR = "demo-workshops/tcw-oracle/terraform-blueprints"
                    }
                    stages{
                        stage("Oracle: Build Unstaged Integrated Test Environment"){
                            when {
                                expression { return (params.ORACLE_TESTING == true && ORACLE_CHANGE == "true") || params.ORACLE_FORCE_TESTING == true }
                            }
                            steps{
                                echo "Copying firewall rules"
                                sh "cp /var/lib/jenkins/firewall.tf ${env.TF_DIR}/modules/firewall/firewall.tf"
                                dir(env.TF_DIR){
                                    script{environment.terraformBuild("false",env.WORKSHOP)}
                                }
                            }
                        }
                        stage('Oracle: Unstaged Integration Testing'){
                            when {
                                expression { return (params.ORACLE_TESTING == true && ORACLE_CHANGE == "true") || params.ORACLE_FORCE_TESTING == true }
                            }
                            steps{
                                dir(env.TF_DIR){
                                    script{environment.environmentReady(env.WORKSHOP)}
                                }
                            }
                        }
                        stage('Oracle: Create Staged AMIs'){
                            when {
                                expression { return (params.ORACLE_TESTING == true && ORACLE_CHANGE == "true") || params.ORACLE_FORCE_TESTING == true }
                            }
                            steps{
                                script{
                                    environment.amiify(env.WORKSHOP)
                                }
                            }
                        }
                        stage('Oracle: Build Staged Integrated Test Environment'){
                            when {
                                expression { return (params.ORACLE_TESTING == true && ORACLE_CHANGE == "true") || params.ORACLE_FORCE_TESTING == true }
                            }
                            steps{
                                dir(env.TF_DIR){
                                    // Change the subnet so that we are also testing the reconfiguration on provisioning
                                    sh "sed -i -e 's|\"0\".*|\"0\" = \"10.0.2.0/24\"|' modules/subnet/variables.tf"
                                    script{environment.terraformBuild("true",env.WORKSHOP)}
                                }
                            }
                        }
                        stage('Oracle: Staged Integration Testing'){
                            when {
                                expression { return (params.ORACLE_TESTING == true && ORACLE_CHANGE == "true") || params.ORACLE_FORCE_TESTING == true }
                            }
                            steps{
                                dir(env.TF_DIR){
                                    script{environment.environmentReady(env.WORKSHOP)}
                                    script{environment.environmentTest(env.WORKSHOP)}
                                }
                            }
                        }
                    }
                }
                stage('Postgres'){
                    environment {
                        WORKSHOP = "tcw-pg"
                        TF_DIR = "demo-workshops/tcw-pg/terraform-blueprints"
                    }
                    stages{
                        stage("Postgres: Build Unstaged Integrated Test Environment"){
                            when {
                                expression { return (params.PG_TESTING == true && PG_CHANGE == "true") || params.PG_FORCE_TESTING == true }
                            }
                            steps{
                                echo "Copying firewall rules"
                                sh "cp /var/lib/jenkins/firewall.tf ${env.TF_DIR}/modules/firewall/main.tf"
                                dir(env.TF_DIR){
                                    script{environment.terraformBuild("false","tcw-pg")}
                                }
                            }
                        }
                        stage('Postgres: Unstaged Integration Testing'){
                            when {
                                expression { return (params.PG_TESTING == true && PG_CHANGE == "true") || params.PG_FORCE_TESTING == true }
                            }
                            steps{
                                dir(env.TF_DIR){
                                    script{environment.environmentReady(env.WORKSHOP)}
                                }
                            }
                        }
                        stage('Postgres: Create Staged AMIs'){
                            when {
                                expression { return (params.PG_TESTING == true && PG_CHANGE == "true") || params.PG_FORCE_TESTING == true }
                            }
                            steps{
                                script{
                                    environment.amiify(env.WORKSHOP)
                                }
                            }
                        }
                        stage('Postgres: Build Staged Integrated Test Environment'){
                            when {
                                expression { return (params.PG_TESTING == true && PG_CHANGE == "true") || params.PG_FORCE_TESTING == true }
                            }
                            steps{
                                dir(env.TF_DIR){
                                    // Change the subnet so that we are also testing the reconfiguration on provisioning
                                    sh "sed -i -e 's|\"0\".*|\"0\" = \"10.0.2.0/24\"|' modules/subnet/variables.tf"
                                    script{environment.terraformBuild("true",env.WORKSHOP)}
                                }
                            }
                        }
                        stage('Postgres: Staged Integration Testing'){
                            when {
                                expression { return (params.PG_TESTING == true && PG_CHANGE == "true") || params.PG_FORCE_TESTING == true }
                            }
                            steps{
                                dir(env.TF_DIR){
                                    script{environment.environmentReady(env.WORKSHOP)}
                                    script{environment.environmentTest(env.WORKSHOP)}
                                }
                            }
                        }
                    }
                }
            }   
        }
    }
    post{
        always{
            script{environment.chownership()}
        }
        success{
            script{environment.terraformDestroy("tcw-oracle")}
            script{environment.terraformDestroy("tcw-pg")}
        }
    }
}
 