pipelineJob('sharran-job') {
  definition {
  cpsScm {
      lightweight()
      scm {
        git('https://github.com/msharran/labs.git', "main", {node -> node / 'extensions' << '' })
      }
      scriptPath("k8s/jenkins-operator/Jenkinsfile.groovy")
    }
  }
}
