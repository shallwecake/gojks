package test

import (
	"fmt"
	"github.com/shallwecake/gojks/jenkins_suggest"
	"testing"
)

var (
	taskName = "declare"
)

func TestInnit(t *testing.T) {

	jenkinsURL := "http://localhost:8500"
	jobName := "test-jenkins-Pipeline"
	username := "admin"
	apiToken := "admin"

	auth := &jenkins_suggest.Auth{
		Username: username,
		ApiToken: apiToken,
	}
	jenkins := jenkins_suggest.NewJenkins(auth, jenkinsURL)
	//jenkins.GetJob()
	names, _ := jenkins.Query(jobName)

	fmt.Println(names)
}
