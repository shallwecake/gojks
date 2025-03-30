package test

import (
	"fmt"
	"github.com/shallwecake/gojks/ifunction"
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

	auth := &ifunction.Auth{
		Username: username,
		ApiToken: apiToken,
	}
	jenkins := ifunction.NewJenkins(auth, jenkinsURL)
	//jenkins.GetJob()
	names, _ := jenkins.Query(jobName)

	fmt.Println(names)
}
