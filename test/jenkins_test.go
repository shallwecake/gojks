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
	auth := &jenkins.Auth{
		Username: "pangwangbin",
		ApiToken: "wongbin123",
	}
	jenkins := jenkins.NewJenkins(auth, "https://jenkins.gw-greenenergy.com")
	//jenkins.GetJob()
	names, _ := jenkins.FuzzyJobName("decl")

	fmt.Println(names)
}
