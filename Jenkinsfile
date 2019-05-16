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
                echo "Installing Ansible role requirements"
                sh "ansible-galaxy install -r demo-workshops/ansible/all_requirements.yml"
                echo "Grabbing certs"
                sh "cp ${env.ANSIBLE_CERT}* certs"
            }
        }
        stage('Stage 1 Packer Builds'){
            parallel {
                stage('delphix-centos7-ansible-base.json'){
                    when {
                        expression {return environment.ami_checker(false) == 0}
                    }
                    steps{
                        script{
                            environment.packerBuild()
                            CHANGE = true
                            }
                    }
                }
                // stage('delphix-ubuntu-bionic-guacamole.json'){
                //     when {
                //         expression {return environment.ami_checker(false) == 0}
                //     }
                //     steps{
                //         script{
                //             environment.packerBuild()
                //             CHANGE = true
                //         }  
                //     }
                // }
            }
        }
        // stage('Stage 2 Packer Builds'){
        //     parallel {
        //         stage('delphix-centos7-daf-app.json'){
        //             when {
        //                 expression {return environment.ami_checker(false) == 0}
        //             }
        //             steps{
        //                 script{
        //                     environment.packerBuild()
        //                     CHANGE = true
        //                 }  
        //             }
        //         }
        //         stage('delphix-centos7-oracle-12.2.0.1.json'){
        //             when {
        //                 expression {return environment.ami_checker(false) == 0}
        //             }
        //             steps{
        //                 script{
        //                     environment.packerBuild()
        //                     CHANGE = true
        //                 }  
        //             }
        //         }
        //         stage('delphix-centos7-kitchen_sink.json'){
        //             when {
        //                 expression {return environment.ami_checker(false) == 0}
        //             }
        //             steps{
        //                 script{
        //                     environment.packerBuild()
        //                     CHANGE = true
        //                 }  
        //             }
        //         }
        //         stage('delphix-tcw-jumpbox.json'){
        //             when {
        //                 expression {return environment.ami_checker(false) == 0}
        //             }
        //             steps{
        //                 script{
        //                     environment.packerBuild()
        //                     CHANGE = true
        //                 }  
        //             }
        //         }
        //     }
        // }
        // stage('Stage 3 Packer Builds'){
        //     parallel {
        //         stage('delphix-tcw-oracle12-source.json'){
        //             when {
        //                 expression {return environment.ami_checker(false) == 0}
        //             }
        //             steps{
        //                 script{
        //                     environment.packerBuild()
        //                     CHANGE = true
        //                 }  
        //             }
        //         }
        //         stage('delphix-tcw-oracle12-target.json'){
        //             when {
        //                 expression {return environment.ami_checker(false) == 0}
        //             }
        //             steps{
        //                 script{
        //                     environment.packerBuild()
        //                     CHANGE = true
        //                 }  
        //             }
        //         }
        //         stage('delphix-centos7-tooling-base.json'){
        //             when {
        //                 expression {return environment.ami_checker(false) == 0}
        //             }
        //             steps{
        //                 script{
        //                     environment.packerBuild()
        //                     CHANGE = true
        //                 }  
        //             }
        //         }
        //     }
        // }
        // stage('Stage 4 Packer Builds'){
        //     parallel {
        //         stage('delphix-tcw-tooling-oracle.json'){
        //             when {
        //                 expression {return environment.ami_checker(false) == 0}
        //             }
        //             steps{
        //                 script{
        //                     environment.packerBuild()
        //                     CHANGE = true
        //                 }  
        //             }
        //         }
        //     }
        // }
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
        success{
            dir(env.TF_DIR){
                script{environment.terraformDestroy()}
            }
        }
    }
}
 