pipelineJob('sharran-job') {
  definition {
  cpsScm {
      lightweight()
      scm {
        git('https://github.com/msharran/labs.git', "main")
        // gitSCM {
        //   userRemoteConfigs {
        //     userRemoteConfig {
        //       url("")
        //       name("labs")
        //       refspec("https://github.com/msharran/labs.git")
        //       credentialsIdj("https://github.com/msharran/labs.git")
        //     }
        //   }
        //   branches {
        //     branchSpec {
        //       name("main")
        //     }
        //   }
        // }
      }
      scriptPath("k8s/jenkins-operator/Jenkinsfile.groovy")
    }
  }
}
