pipelineJob('sharran-job') {
  definition {
  cpsScm {
      lightweight()
      scm {
        gitSCM {
          userRemoteConfigs {
            userRemoteConfig {
              url("https://github.com/msharran/labs.git")
            }
          }
          branches {
            branchSpec {
              name("main")
            }
          }
        }
      }
      scriptPath("k8s/jenkins-operator/Jenkinsfile.groovy")
    }
  }
}
