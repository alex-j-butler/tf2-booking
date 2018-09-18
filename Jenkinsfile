def build(os, arch) {
	withEnv(["GOOS=${os}", "GOARCH=${arch}"]) {
		sh """cd $GOPATH/src/alex-j-butler/tf2-booking/ && go build -o tf2-booking-$os-$arch"""
		archiveArtifacts artifacts: "tf2-booking-$os-$arch", fingerprint: true
	}
}

node {
	try {
		def root = tool name: 'Go1.10.3', type: 'go'
 
		// Export environment variables pointing to the directory where Go was installed
		withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {

				ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/alex-j-butler/tf2-booking/") {
					withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
						env.PATH="${GOPATH}/bin:$PATH"

						stage('Checkout') {
							checkout scm
						}

						stage('Dependencies') {
							sh """cd $GOPATH/src/alex-j-butler/tf2-booking/ && go get github.com/kardianos/govendor && govendor sync"""
						}

					stage('Build') {
						parallel (
							linuxamd64: {
								build('linux', 'amd64')
							}
						)
					}
				}
			}
		} catch (e) {
			currentBuild.result = "FAILED"
		}
	}
}