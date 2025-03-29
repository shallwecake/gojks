package test

import (
	"fmt"
	"gojks/jenkins"
	"testing"
)

var (
	taskName = "declare"
)

// jenkins, err := gojenkins.CreateJenkins(nil, "https://jenkins.gw-greenenergy.com/", "pangwangbin", "wongbin123").Init(ctx)
func TestInnit(t *testing.T) {

	jenkinsURL := "http://localhost:8500"
	jobName := "test-jenkins-Pipeline"
	username := "admin"
	apiToken := "admin"

	auth := &jenkins.Auth{
		Username: username,
		ApiToken: apiToken,
	}
	jenkins := jenkins.NewJenkins(auth, jenkinsURL)
	//jenkins.GetJob()
	names, _ := jenkins.Query(jobName)

	fmt.Println(names)
}
