@Library('environment-lib') _

pipeline {
    options {
        disableConcurrentBuilds()
    }
    agent any
    environment {
        // ANSIBLE_ROLES_PATH = "ansible/roles"
        CHANGE = false
        TF_DIR = "demo-workshops/tcw/terraform-blueprints"
        PACKER_BASE = "base-templates"
        PACKER_TCW = "demo-workshops/tcw/packer-templates"
        IMAGEBUILDER_LIB = "demo-workshops/tcw/image-builder/library.sh"
        TF_VARS = "/var/lib/jenkins/terraform.tfvars"
        ANSIBLE_CERT = "/var/lib/jenkins/.ssh/ansible"
    }
    parameters {
        booleanParam(name: 'TESTING', defaultValue: true, description: 'Run Testing on Change')
        booleanParam(name: 'FORCE_TESTING', defaultValue: false, description: 'Force Testing')
    }
    stages {
        stage('Prep Environment'){
            steps {
                // echo "Installing Ansible role requirements"
                // sh "ansible-galaxy install -r demo-workshops/ansible/all_requirements.yml"
                echo "Grabbing certs"
                sh "cp ${env.ANSIBLE_CERT}* certs"
                echo "Copying env file"
                sh "cp /var/lib/jenkins/.docker.env ."
                echo "Copying s3 backend file"
                sh "cp /var/lib/jenkins/terraform_dev_backend.tf ${env.TF_DIR}/."
                echo "Copying firewall rules"
                sh "cp /var/lib/jenkins/firewall.tf ${env.TF_DIR}/modules/firewall/main.tf"
                echo "Building Container"
                sh "docker-compose build --no-cache" 
            }
        }
        stage('Packer Builds'){
            steps{
                script{
                    environment.packerBuild()
                    CHANGE = sh (
                        script: """#!/bin/bash
                                { set +x; } 2>/dev/null
                                [[ -f change.ignore ]] && echo "true" || echo "false" """,
                        returnStdout: true
                    ).trim()
                    }
            }
        }
        stage('Build Unstaged Integrated Test Environment'){
            when {
                expression { return (params.TESTING == true && CHANGE == true) || params.FORCE_TESTING == true }
            }
            steps{
                dir(env.TF_DIR){
                    script{environment.terraformBuild("false")}
                }
            }
        }
        stage('Unstaged Integration Testing'){
            when {
                expression { return (params.TESTING == true && CHANGE == true) || params.FORCE_TESTING == true }
            }
            steps{
                dir(env.TF_DIR){
                    script{environment.environmentTest()}
                }
            }
        }
        stage('Create Staged AMIs'){
            when {
                expression { return (params.TESTING == true && CHANGE == true) || params.FORCE_TESTING == true }
            }
            parallel{
                stage('delphix-tcw-jumpbox'){
                    steps{
                        dir(env.TF_DIR){
                            script{environment.amiify()}
                        }
                    }
                }
                stage('delphix-tcw-oracle12-source'){
                    steps{
                        dir(env.TF_DIR){
                            script{environment.amiify()}
                        }
                    }
                }
                stage('delphix-tcw-oracle12-target'){
                    steps{
                        dir(env.TF_DIR){
                            script{environment.amiify()}
                        }
                    }
                }
                stage('delphix-tcw-tooling-oracle'){
                    steps{
                        dir(env.TF_DIR){
                            script{environment.amiify()}
                        }
                    }
                }
                stage('delphix-tcw-delphixengine'){
                    steps{
                        dir(env.TF_DIR){
                            script{environment.amiify()}
                        }
                    }
                }
            }
        }
        stage('Build Staged Integrated Test Environment'){
            when {
                expression { return (params.TESTING == true && CHANGE == true) || params.FORCE_TESTING == true }
            }
            steps{
                dir(env.TF_DIR){
                    script{environment.terraformBuild("true")}
                }
            }
        }
        stage('Staged Integration Testing'){
            when {
                expression { return (params.TESTING == true && CHANGE == true) || params.FORCE_TESTING == true }
            }
            steps{
                dir(env.TF_DIR){
                    script{environment.environmentTest()}
                }
            }
        }
    }
    post{
        always{
            script{environment.terraformDestroy()}
        }
    }
}
 