package test

import (
	"github.com/shallwecake/gojks/ifunction"
	"testing"
)

func TestRancher01(t *testing.T) {
	engine := ifunction.InitDb()
	defer ifunction.CloseDbEngine(engine)

	ifunction.MonitorRancher("test")
}
