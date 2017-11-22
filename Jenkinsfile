node('master') {
  stage('Unit Tests') {
    git url: "https://github.com/bobbydeveaux/micro-user-rest.git"
  }
  stage('Build Bin') {
    sh "go get -v -d ./..."
    sh "CGO_ENABLED=0 GOOS=linux go build -o micro-user-rest ."
  }
  stage('Build Image') {
    sh "oc start-build micro-user-rest --from-file=. --follow"
  }
  stage('Deploy') {
    openshiftDeploy depCfg: 'micro-user-rest', namespace: 'fbac'
    openshiftVerifyDeployment depCfg: 'micro-user-rest', replicaCount: 1, verifyReplicaCount: true
  }
}